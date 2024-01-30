package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm/clause"
	"os/exec"
	"strings"
	"tiktok/src/constant/config"
	strings2 "tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/models"
	database "tiktok/src/storage/db"
	"tiktok/src/storage/file"
	"tiktok/src/utils/logging"
	"tiktok/src/utils/pathgen"
	"tiktok/src/utils/rabbitmq"
)

var openaiClient = openai.NewClient(config.EnvCfg.ChatGPTAPIKEYS)
var delayTime = int32(2 * 60 * 1000) //2 minutes
var maxRetries = int32(3)

// errorHandler If `requeue` is false, it will just `Nack` it. If `requeue` is true, it will try to re-publish it.
func errorHandler(channel *amqp.Channel, d amqp.Delivery, requeue bool, logger *logrus.Entry, span *trace.Span) {
	if !requeue { // Nack the message
		err := d.Nack(false, false)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error when nacking the video...")
			logging.SetSpanError(*span, err)
		}
	} else { // Re-publish the message
		curRetry, ok := d.Headers["x-retry"].(int32)
		if !ok {
			curRetry = 0
		}
		if curRetry >= maxRetries {
			logger.WithFields(logrus.Fields{
				"body": d.Body,
			}).Errorf("Maximum retries reached for message.")
			logging.SetSpanError(*span, errors.New("maximum retries reached for message"))
			err := d.Ack(false)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when dealing with the video...")
			}
		} else {
			curRetry++
			headers := d.Headers
			headers["x-delay"] = delayTime
			headers["x-retry"] = curRetry

			err := d.Ack(false)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when dealing with the video...")
			}

			logger.Debugf("Retrying %d times", curRetry)

			err = channel.Publish(
				strings2.VideoExchange,
				strings2.VideoSummary,
				false,
				false,
				amqp.Publishing{
					DeliveryMode: amqp.Persistent,
					ContentType:  "text/plain",
					Body:         d.Body,
					Headers:      headers,
				},
			)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err": err,
				}).Errorf("Error when re-publishing the video to queue...")
				logging.SetSpanError(*span, err)
			}
		}
	}
}

func SummaryConsume(channel *amqp.Channel) {
	msg, err := channel.Consume(strings2.VideoSummary, "", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	for d := range msg {
		//解包 Otel Context
		ctx := rabbitmq.ExtractAMQPHeaders(context.Background(), d.Headers)
		ctx, span := tracing.Tracer.Start(ctx, "VideoSummaryService")
		logger := logging.LogService("VideoSummary").WithContext(ctx)

		var raw models.RawVideo
		if err := json.Unmarshal(d.Body, &raw); err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error when unmarshaling the prepare json body.")
			logging.SetSpanError(span, err)

			errorHandler(channel, d, false, logger, &span)
			span.End()
			continue
		}
		logger.WithFields(logrus.Fields{
			"RawVideo": raw.VideoId,
		}).Debugf("Receive message of video %d", raw.VideoId)

		// Video -> Audio
		audioFileName := pathgen.GenerateAudioName(raw.FileName)
		isAudioFileExist, _ := file.IsFileExist(ctx, audioFileName)
		if !isAudioFileExist {
			audioFileName, err = video2Audio(ctx, raw.FileName)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err":             err,
					"video_file_name": raw.FileName,
				}).Errorf("Failed to transform video to audio")
				logging.SetSpanError(span, err)

				errorHandler(channel, d, false, logger, &span)
				span.End()
				continue
			}
		} else {
			logger.WithFields(logrus.Fields{
				"VideoId": raw.VideoId,
			}).Debugf("Video %d already exists audio", raw.VideoId)
		}

		// Audio -> Transcript
		transcriptExist, transcript, err := isTranscriptExist(raw.VideoId)
		if err != nil {
			logger.WithFields(
				logrus.Fields{
					"err":     err,
					"VideoId": raw.VideoId,
				}).Errorf("Faild to get transcript of video %d from database", raw.VideoId)
		}
		if !transcriptExist {
			transcript, err = speech2Text(ctx, audioFileName)
			if err != nil {
				logger.WithFields(logrus.Fields{
					"err":             err,
					"audio_file_name": audioFileName,
				}).Errorf("Failed to get transcript of an audio from ChatGPT")
				logging.SetSpanError(span, err)

				errorHandler(channel, d, true, logger, &span)
				span.End()
				continue
			}
		} else {
			logger.WithFields(logrus.Fields{
				"VideoId": raw.VideoId,
			}).Debugf("Video %d already exists transcript", raw.VideoId)
		}

		var (
			summary  string
			keywords string
		)

		// Transcript -> Summary
		summaryChannel := make(chan string)
		summaryErrChannel := make(chan error)
		summaryExist, summary, err := isSummaryExist(raw.VideoId)
		if err != nil {
			logger.WithFields(
				logrus.Fields{
					"err":     err,
					"VideoId": raw.VideoId,
				}).Errorf("Faild to get summary of video %d from database", raw.VideoId)
		}
		if !summaryExist {
			go text2Summary(ctx, transcript, &summaryChannel, &summaryErrChannel)
		} else {
			logger.WithFields(logrus.Fields{
				"VideoId": raw.VideoId,
			}).Debugf("Video %d already exists summary", raw.VideoId)
		}

		// Transcript -> Keywords
		keywordsChannel := make(chan string)
		keywordsErrChannel := make(chan error)
		keywordsExist, keywords, err := isKeywordsExist(raw.VideoId)
		if err != nil {
			logger.WithFields(
				logrus.Fields{
					"err":     err,
					"VideoId": raw.VideoId,
				}).Errorf("Faild to get keywords of video %d from database", raw.VideoId)
		}
		if !keywordsExist {
			go text2Keywords(ctx, transcript, &keywordsChannel, &keywordsErrChannel)
		} else {
			logger.WithFields(logrus.Fields{
				"VideoId": raw.VideoId,
			}).Debugf("Video %d already exists keywords", raw.VideoId)
		}

		summaryOrKeywordsErr := false

		if !summaryExist {
			select {
			case summary = <-summaryChannel:
			case err = <-summaryErrChannel:
				logger.WithFields(logrus.Fields{
					"err":             err,
					"audio_file_name": audioFileName,
				}).Errorf("Failed to get summary of an audio from ChatGPT")
				logging.SetSpanError(span, err)
				summary = ""
				summaryOrKeywordsErr = true
			}
		}

		if !keywordsExist {
			select {
			case keywords = <-keywordsChannel:
			case err = <-keywordsErrChannel:
				logger.WithFields(logrus.Fields{
					"err":             err,
					"audio_file_name": audioFileName,
				}).Errorf("Failed to get keywords of an audio from ChatGPT")
				logging.SetSpanError(span, err)
				keywords = ""
				summaryOrKeywordsErr = true
			}
		}

		// Update summary information to database
		video := &models.Video{
			ID:            raw.VideoId,
			AudioFileName: audioFileName,
			Transcript:    transcript,
			Summary:       summary,
			Keywords:      keywords,
		}
		result := database.Client.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"audio_file_name", "transcript", "summary", "keywords"}),
		}).Create(&video)
		if result.Error != nil {
			logger.WithFields(logrus.Fields{
				"Err":           result.Error,
				"ID":            raw.VideoId,
				"AudioFileName": audioFileName,
				"Transcript":    transcript,
				"Summary":       summary,
				"Keywords":      keywords,
			}).Errorf("Error when updating summary information to database")
			logging.SetSpanError(span, result.Error)
			errorHandler(channel, d, true, logger, &span)
			span.End()
			continue
		}

		// Cannot get summary or keywords from ChatGPT: resend to mq
		if summaryOrKeywordsErr {
			errorHandler(channel, d, true, logger, &span)
			span.End()
			continue
		}

		span.End()
		err = d.Ack(false)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error when dealing with the video...")
		}
	}
}

func video2Audio(ctx context.Context, videoFileName string) (audioFileName string, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "Video2Audio")
	defer span.End()
	logger := logging.LogService("VideoSummary.Video2Audio").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"video_file_name": videoFileName,
	}).Debugf("Transforming video to audio")

	videoFilePath := file.GetLocalPath(ctx, videoFileName)
	cmdArgs := []string{
		"-i", videoFilePath, "-q:a", "0", "-map", "a", "-f", "mp3", "-",
	}
	cmd := exec.Command("ffmpeg", cmdArgs...)
	var buf bytes.Buffer
	cmd.Stdout = &buf

	err = cmd.Run()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"VideoFileName": videoFileName,
		}).Errorf("cmd %s failed with %s", "ffmpeg "+strings.Join(cmdArgs, " "), err)
		logging.SetSpanError(span, err)
		return
	}

	audioFileName = pathgen.GenerateAudioName(videoFileName)

	_, err = file.Upload(ctx, audioFileName, bytes.NewReader(buf.Bytes()))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"VideoFileName": videoFileName,
			"AudioFileName": audioFileName,
		}).Errorf("Failed to upload audio file")
		logging.SetSpanError(span, err)
		return
	}
	return
}

func speech2Text(ctx context.Context, audioFileName string) (transcript string, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "Speech2Text")
	defer span.End()
	logger := logging.LogService("VideoSummary.Speech2Text").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"AudioFileName": audioFileName,
	}).Debugf("Transforming audio to transcirpt")

	audioFilePath := file.GetLocalPath(ctx, audioFileName)

	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: audioFilePath,
	}
	resp, err := openaiClient.CreateTranscription(ctx, req)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Err":           err,
			"AudioFileName": audioFileName,
		}).Errorf("Failed to get transcript from ChatGPT")
		logging.SetSpanError(span, err)
		return
	}

	transcript = resp.Text
	logger.WithFields(logrus.Fields{
		"Transcript": transcript,
	}).Debugf("Successful to get transcript from ChatGPT")

	return
}

func text2Summary(ctx context.Context, transcript string, summaryChannel *chan string, errChannel *chan error) {
	ctx, span := tracing.Tracer.Start(ctx, "Text2Summary")
	defer span.End()
	logger := logging.LogService("VideoSummary.Text2Summary").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"transcript": transcript,
	}).Debugf("Getting transcript summary form ChatGPT")

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You will be provided with a block of text which is the content of a video, " +
					"and your task is to give 2 Simplified Chinese sentences to summarize the video.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: transcript,
			},
		},
	}
	resp, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Err":        err,
			"Transcript": transcript,
		}).Errorf("Failed to get summary from ChatGPT")
		logging.SetSpanError(span, err)
		*errChannel <- err
		return
	}

	summary := resp.Choices[0].Message.Content
	*summaryChannel <- summary

	logger.WithFields(logrus.Fields{
		"Summary": summary,
	}).Debugf("Successful to get summary from ChatGPT")
}

func text2Keywords(ctx context.Context, transcript string, keywordsChannel *chan string, errChannel *chan error) {
	ctx, span := tracing.Tracer.Start(ctx, "Text2Keywords")
	defer span.End()
	logger := logging.LogService("VideoSummary.Text2Keywords").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"transcript": transcript,
	}).Debugf("Getting transcript keywords from ChatGPT")

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: "You will be provided with a block of text which is the content of a video, " +
					"and your task is to give 5 tags in Simplified Chinese to the video to attract audience. " +
					"For example, 美食 | 旅行 | 阅读",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: transcript,
			},
		},
	}
	resp, err := openaiClient.CreateChatCompletion(ctx, req)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Err":        err,
			"Transcript": transcript,
		}).Errorf("Failed to get keywords from ChatGPT")
		logging.SetSpanError(span, err)
		*errChannel <- err
		return
	}

	keywords := resp.Choices[0].Message.Content

	*keywordsChannel <- keywords

	logger.WithFields(logrus.Fields{
		"Keywords": keywords,
	}).Debugf("Successful to get keywords from ChatGPT")
}

func isTranscriptExist(videoId uint32) (res bool, transcript string, err error) {
	video := &models.Video{}
	result := database.Client.Select("transcript").Where("id = ?", videoId).Find(video)
	err = result.Error
	if result.Error != nil {
		res = false
	} else {
		res = video.Transcript != ""
		transcript = video.Transcript
	}
	return
}

func isSummaryExist(videoId uint32) (res bool, summary string, err error) {
	video := &models.Video{}
	result := database.Client.Select("summary").Where("id = ?", videoId).Find(video)
	err = result.Error
	if result.Error != nil {
		res = false
	} else {
		res = video.Summary != ""
		summary = video.Summary
	}
	return
}

func isKeywordsExist(videoId uint32) (res bool, keywords string, err error) {
	video := &models.Video{}
	result := database.Client.Select("keywords").Where("id = ?", videoId).Find(video)
	err = result.Error
	if result.Error != nil {
		res = false
	} else {
		res = video.Keywords != ""
		keywords = video.Keywords
	}
	return
}
