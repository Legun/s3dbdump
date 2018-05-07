package main

import (
    "fmt"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/aws/awsutil"
    "os/exec"
    "log"
    "io/ioutil"
    "time"
)

func s3_input(bucketName string){

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-west-2")},)

    svc := s3.New(sess)

    //Ищем список файлов в папке с дампом
    files, err := ioutil.ReadDir("./DUMP")
    if err != nil {
        log.Fatal(err)
    }

    for _, file := range files {
        fmt.Println(file.Name())

        filename := file.Name()

        //создаем запрос в s3
        input := &s3.PutObjectInput{
            Bucket: aws.String(bucketName),
            Key:    aws.String(filename),
            //Body:   filename,
        }

        resp, err := svc.PutObject(input)
        if err != nil {
            fmt.Printf("bad response: %s", err)
        }

        fmt.Printf("response %s", awsutil.StringValue(resp))
    }
}

func dump_db(){
    // вызываем команду из SH для дампа базы и упаковываем его в архив
    cmd := exec.Command("sh", "-c", "mongodump --db test --gzip --archive=./DUMP/`date +'%m-%d-%y'.gz` ")
    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", stdoutStderr)
}




func main() {

    //Инициализируем бесконечный цыкл
    for {
        //Задаем в качестве аргумента для  NewTimer время ожидания
        // time.Minute * 2 - ожидать 2 минуты
        // time.Hour * 1 - ожидать 1 час
        // time.Hour * 24 ожидать сутки

        timer := time.NewTimer(time.Minute * 2)
        <-timer.C
        dump_db()                 // делаем дамп базы
        s3_input("threew82-demo") // отправляем дамп в хранилище "bucketName" на s3
    }


}