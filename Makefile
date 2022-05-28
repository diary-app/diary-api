heroku-push:
	heroku container:push web -a egor-diary-api
heroku-release:
	heroku container:release web -a egor-diary-api
heroku-deploy:
	make heroku-push
	make heroku-release