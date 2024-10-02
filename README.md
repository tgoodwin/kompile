# Kompile

Try it out..
```
go run main.go servicegen.go --file examples/photoservice/main.go
```

calling the example photoservice:
```
curl -X POST http://localhost:8080/upload -F "image=@/Users/$USER/Desktop/somepicture.jpg"
```
