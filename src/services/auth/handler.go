package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"strconv"
	stringsLib "strings"
	"sync"
	"tiktok/src/constant/strings"
	"tiktok/src/extra/tracing"
	"tiktok/src/models"
	"tiktok/src/rpc/auth"
	"tiktok/src/storage/cached"
	"tiktok/src/storage/db"
	"tiktok/src/utils/logging"
)

type AuthServiceImpl struct {
	auth.AuthServiceServer
}

// Authenticate 方法用于验证用户的token是否有效
func (a AuthServiceImpl) Authenticate(ctx context.Context, request *auth.AuthenticateRequest) (resp *auth.AuthenticateResponse, err error) {
	// 开始一个新的追踪span
	ctx, span := tracing.Tracer.Start(ctx, "AuthenticateService")

	// 确保span在函数结束时关闭
	defer span.End()
	logger := logging.LogService("AuthService.Authenticate").WithContext(ctx)
	// 检查token是否存在
	userId, ok := hasToken(ctx, request.Token)

	if !ok {
		resp = &auth.AuthenticateResponse{
			StatusCode: strings.UserNotExistedCode,
			StatusMsg:  strings.UserNotExisted,
		}
		return
	}

	// 尝试将用户ID从字符串转换为无符号32位整数
	id, err := strconv.ParseUint(userId, 10, 32)
	// 如果在转换过程中发生错误，返回内部错误的响应
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":   err,
			"token": request.Token,
		}).Warnf("AuthService Authenticate Action failed to response when parsering uint")
		span.RecordError(err)
		resp = &auth.AuthenticateResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	// 创建一个新的响应，设置状态码为OK的状态码，状态消息为OK的状态消息，用户ID为转换后的用户ID
	resp = &auth.AuthenticateResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		UserId:     uint32(id),
	}

	// 返回响应
	return
}

// Register 方法用于注册新用户
func (a AuthServiceImpl) Register(ctx context.Context, request *auth.RegisterRequest) (resp *auth.RegisterResponse, err error) {
	// 开始一个新的追踪span
	ctx, span := tracing.Tracer.Start(ctx, "RegisterService")
	// 确保span在函数结束时关闭
	defer span.End()
	logger := logging.LogService("AuthService.Register").WithContext(ctx)

	// 初始化响应
	resp = &auth.RegisterResponse{}
	var user models.User
	// 在数据库中查找是否已存在同名用户
	result := database.Client.WithContext(ctx).Limit(1).Where("user_name = ?", request.Username).Find(&user)
	// 如果找到了同名用户，返回用户已存在的响应
	if result.RowsAffected != 0 {
		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthUserExistedCode,
			StatusMsg:  strings.AuthUserExisted,
		}
		return
	}

	// 对用户密码进行哈希处理
	var hashedPassword string
	if hashedPassword, err = hashPassword(ctx, request.Password); err != nil {
		// 如果在哈希处理过程中发生错误，返回内部错误的响应
		logger.WithFields(logrus.Fields{
			"err":      result.Error,
			"username": request.Username,
		}).Warnf("AuthService Register Action failed to response when hashing password")
		span.RecordError(err)
		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	// 使用WaitGroup来同步并发的goroutine
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 在一个新的goroutine中获取用户签名
	go func() {
		defer wg.Done()
		// 从hitokoto服务获取一句话作为用户签名
		resp, err := http.Get("https://v1.hitokoto.cn/?c=b&encode=text")
		_, span := tracing.Tracer.Start(ctx, "FetchSignature")
		defer span.End()
		logger := logging.LogService("AuthService.FetchSignature").WithContext(ctx)

		if err != nil {
			user.Signature = user.UserName
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("Can not reach hitokoto")
			span.RecordError(err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			user.Signature = user.UserName
			logger.WithFields(logrus.Fields{
				"status_code": resp.StatusCode,
			}).Warnf("Hitokoto service may be error")
			span.RecordError(err)
			return
		}

		// 读取hitokoto服务返回的响应体作为用户签名
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			// 如果在读取响应体过程中发生错误，使用用户名作为签名
			user.Signature = user.UserName
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Warnf("Can not decode the response body of hitokoto")
			span.RecordError(err)
			return
		}

		// 设置用户签名
		user.Signature = string(body)
	}()

	// 在一个新的goroutine中获取用户头像
	go func() {
		defer wg.Done()
		// 如果用户名是一个邮箱地址，使用邮箱地址获取用户头像
		if user.IsNameEmail() {
			user.Avatar = getAvatarByEmail(ctx, user.UserName)
		}
	}()

	// 等待所有goroutine完成
	wg.Wait()

	// 设置用户名和密码
	user.UserName = request.Username
	user.Password = hashedPassword

	// 在数据库中创建新用户
	result = database.Client.WithContext(ctx).Create(&user)
	if result.Error != nil {
		// 如果在创建用户过程中发生错误，返回内部错误的响应
		logger.WithFields(logrus.Fields{
			"err":      result.Error,
			"username": request.Username,
		}).Warnf("AuthService Register Action failed to response when creating user")
		span.RecordError(err)
		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	// 获取用户token
	if resp.Token = getToken(ctx, user.ID); err != nil {
		// 如果在获取token过程中发生错误，返回内部错误的响应
		logger.WithFields(logrus.Fields{
			"err":      result.Error,
			"username": request.Username,
		}).Warnf("AuthService Register Action failed to response when getting token")
		span.RecordError(err)
		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return resp, nil
	}

	// 设置响应的用户ID，状态码和状态消息
	resp.UserId = uint32(user.ID)
	resp.StatusCode = strings.ServiceOKCode
	resp.StatusMsg = strings.ServiceOK
	return
}

// Login 方法用于验证用户的登录信息
func (a AuthServiceImpl) Login(ctx context.Context, request *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	// 开始一个新的追踪span
	ctx, span := tracing.Tracer.Start(ctx, "LoginService")
	// 确保span在函数结束时关闭
	defer span.End()
	logger := logging.LogService("AuthService.Login").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"username": request.Username,
	}).Infof("User try to log in.")

	// 初始化响应
	resp = &auth.LoginResponse{}
	user := models.User{
		UserName: request.Username,
	}
	// 在Redis中验证用户信息
	if !isUserVerifiedInRedis(ctx, request.Username, request.Password) {
		// 在数据库中查找用户
		result := database.Client.Where("user_name = ?", request.Username).WithContext(ctx).Find(&user)

		logger.WithFields(logrus.Fields{
			"err":      result.Error,
			"username": request.Username,
		}).Warnf("AuthService Login Action failed to response with inner err.")
		span.RecordError(result.Error)
		// 如果在查找过程中发生错误，返回内部错误的响应
		if result.Error != nil {
			resp = &auth.LoginResponse{
				StatusCode: strings.AuthServiceInnerErrorCode,
				StatusMsg:  strings.AuthServiceInnerError,
			}
			return
		}

		// 如果没有找到用户，返回用户不存在的响应
		if result.RowsAffected == 0 {
			resp = &auth.LoginResponse{
				StatusCode: strings.UserNotExistedCode,
				StatusMsg:  strings.UserNotExisted,
			}
			return
		}

		// 验证用户密码
		if !checkPasswordHash(ctx, request.Password, user.Password) {
			resp = &auth.LoginResponse{
				StatusCode: strings.AuthUserLoginFailedCode,
				StatusMsg:  strings.AuthUserLoginFailed,
			}
			return
		}

		// 对用户密码进行哈希处理
		hashed, errs := hashPassword(ctx, request.Password)
		if errs != nil {
			logger.WithFields(logrus.Fields{
				"err":      errs,
				"username": request.Username,
			}).Warnf("AuthService Login Action failed to response with inner err.")
			span.RecordError(errs)
			resp = &auth.LoginResponse{
				StatusCode: strings.AuthServiceInnerErrorCode,
				StatusMsg:  strings.AuthServiceInnerError,
			}
			return
		}
		// 将用户信息存储到Redis
		setUserInfoToRedis(ctx, user.UserName, hashed)
	}

	// 获取用户token
	token := getToken(ctx, user.ID)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"username": request.Username,
		}).Warnf("AuthService Login Action failed to response with inner err.")
		span.RecordError(err)
		resp = &auth.LoginResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return resp, nil
	}

	// 设置响应的用户ID，状态码，状态消息和token
	resp = &auth.LoginResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		UserId:     uint32(user.ID),
		Token:      token,
	}
	return
}

func hashPassword(ctx context.Context, password string) (string, error) {
	_, span := tracing.Tracer.Start(ctx, "PasswordHash")
	defer span.End()

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkPasswordHash(ctx context.Context, password, hash string) bool {
	_, span := tracing.Tracer.Start(ctx, "PasswordHashChecked")
	defer span.End()
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getToken(ctx context.Context, userId uint32) string {
	return cached.GetWithFunc(ctx, "U2T"+strconv.FormatUint(uint64(userId), 10),
		func(ctx context.Context, key string) string {
			span := trace.SpanFromContext(ctx)
			token := uuid.New().String()
			span.SetAttributes(attribute.String("token", token))
			cached.Write(ctx, "T2U"+token, strconv.FormatUint(uint64(userId), 10), true)
			return token
		})
}

func hasToken(ctx context.Context, token string) (string, bool) {
	return cached.Get(ctx, "T2U"+token)
}

func isUserVerifiedInRedis(ctx context.Context, username string, password string) bool {
	pass, ok := cached.Get(ctx, "UserLog"+username)
	if !ok {
		return false
	}

	if checkPasswordHash(ctx, password, pass) {
		return true
	}
	return false
}

func setUserInfoToRedis(ctx context.Context, username string, password string) {
	if _, ok := cached.Get(ctx, "UserLog"+username); ok {
		cached.TagDelete(ctx, "UserLog"+username)
	}
	cached.Write(ctx, "UserLog"+username, password, true)
}

func getAvatarByEmail(ctx context.Context, email string) string {
	ctx, span := tracing.Tracer.Start(ctx, "Auth-GetAvatar")
	defer span.End()
	return fmt.Sprintf("https://cravatar.cn/avatar/%s?d=identicon", getEmailMD5(ctx, email))
}

func getEmailMD5(ctx context.Context, email string) (md5String string) {
	_, span := tracing.Tracer.Start(ctx, "Auth-EmailMD5")
	defer span.End()

	lowerEmail := stringsLib.ToLower(email)
	hashed := md5.New()
	hashed.Write([]byte(lowerEmail))
	md5Bytes := hashed.Sum(nil)
	md5String = hex.EncodeToString(md5Bytes)
	return
}
