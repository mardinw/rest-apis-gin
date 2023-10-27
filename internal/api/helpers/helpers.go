package helpers

import (
	"errors"
	"log"

	"payuoge.com/pkg/aws"
)

var ErrPermission = errors.New("permission denied")

func CheckAccountRetail(username *string) error {

	resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*username)
	if err != nil {
		return err
	}

	targetValue := "retail"
	found := false
	for _, value := range resp {
		if value == targetValue {
			found = true
			break
		}
	}

	if !found {
		return ErrPermission
	} else {
		log.Printf("Found %s\n", targetValue)
	}

	return nil

}

func CheckAccountGroceries(username *string) error {

	resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*username)
	if err != nil {
		return err
	}

	targetValue := "grosir"
	found := false
	for _, value := range resp {
		if value == targetValue {
			found = true
			break
		}
	}

	if !found {
		return ErrPermission
	} else {
		log.Printf("Found %s\n", targetValue)
	}

	return nil
}

func CheckAccountAdmin(username *string) error {

	resp, err := aws.NewConnect().Cognito.CheckUserInGroup(*username)
	if err != nil {
		return err
	}

	targetValue := "admin"
	found := false
	for _, value := range resp {
		if value == targetValue {
			found = true
			break
		}
	}

	if !found {
		return ErrPermission
	} else {
		log.Printf("Found %s\n", targetValue)
	}

	return nil
}
