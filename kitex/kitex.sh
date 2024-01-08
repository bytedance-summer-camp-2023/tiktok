kitex -module tiktok -I ./ -v -service usersrv user.proto

kitex -module tiktok -I ./ -v -service commentsrv comment.proto

kitex -module tiktok -I ./ -v -service relationsrv relation.proto

kitex -module tiktok -I ./ -v -service favoritesrv favorite.proto

kitex -module tiktok -I ./ -v -service videosrv video.proto
