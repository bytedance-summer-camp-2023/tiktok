kitex -module github.com/bytedance-summer-camp-2023/tiktok -I ./ -v -service usersrv user.proto

kitex -module github.com/bytedance-summer-camp-2023/tiktok -I ./ -v -service commentsrv comment.proto

kitex -module github.com/bytedance-summer-camp-2023/tiktok -I ./ -v -service relationsrv relation.proto

kitex -module github.com/bytedance-summer-camp-2023/tiktok -I ./ -v -service favoritesrv favorite.proto

kitex -module github.com/bytedance-summer-camp-2023/tiktok -I ./ -v -service videosrv video.proto
