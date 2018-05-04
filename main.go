package main

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"os/exec"
	"log"
	"time"
)

func s3_input(){
	aws_access_key_id := "Input_ID_Key"
	aws_secret_access_key := "Input_Secret_key"
	token := ""
	creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)

	_, err := creds.Get()
	if err != nil {
		fmt.Printf("bad credentials: %s", err)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},)

	svc := s3.New(sess)

	dump_db()

	file, err := os.Open("05-04-18.tar.xz")

	if err != nil {
		fmt.Printf("err opening file: %s", err)
	}

	defer file.Close()


	path := file.Name()

	input := &s3.PutObjectInput{
		Bucket: aws.String("test.threew82"),
		Key:    aws.String(path),
		Body:   file,
	}

	resp, err := svc.PutObject(input)
	if err != nil {
		fmt.Printf("bad response: %s", err)
	}

	fmt.Printf("response %s", awsutil.StringValue(resp))
}

func dump_db(){
	cmd := exec.Command("sh", "-c", "mongodump --db test --gzip --archive=`date +'%m-%d-%y'.gz` ")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}


	func timer(hour, min, sec int, f func()) error {
		loc, err := time.LoadLocation("Local")
		if err != nil {
		return err
	}

		// Вычисляем время первого запуска.
		now := time.Now().Local()
		firstCallTime := time.Date(
		now.Year(), now.Month(), now.Day(), hour, min, sec, 0, loc)
		if firstCallTime.Before(now) {
		// Если получилось время раньше текущего, прибавляем сутки.
		firstCallTime = firstCallTime.Add(time.Hour * 24)
	}

		// Вычисляем временной промежуток до запуска.
		duration := firstCallTime.Sub(time.Now().Local())

		go func() {
		time.Sleep(duration)
		for {
		f()
		// Следующий запуск через сутки.
		time.Sleep(time.Hour * 24)
	}
	}()

		return nil
	}

func main() {
	err := timer(0, 0, 0, s3_input)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Эмуляция дальнейшей работы программы.
	time.Sleep(time.Hour * 24)

}