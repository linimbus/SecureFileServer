curl -X POST -F "file=@readme.txt" -H "Authorization: xxx" http://xxx:8080/upload
curl -X GET http://xxx:8080/download/308bb140-1393-f347-8eaf-e433c4bcb60c -o downloaded_readme.txt