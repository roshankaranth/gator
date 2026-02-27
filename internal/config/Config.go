package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const filename string = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() (Config, error) {
	file_path, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	jsonData, err := os.ReadFile(file_path)

	if err != nil {
		return Config{}, fmt.Errorf("Error while reading file : %v", err)
	}

	var data Config
	err = json.Unmarshal(jsonData, &data)

	if err != nil {
		return Config{}, fmt.Errorf("Error while decoding file : %v", err)
	}

	return data, nil

}

func (c *Config) SetUser(user_name string) error {

	c.Current_user_name = user_name

	err := write(*c)
	if err != nil {
		return err
	}

	return nil

}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("Unable to find home directory : %v", err)
	}

	file_path := home_dir + "/" + filename

	return file_path, nil
}

func write(cfg Config) error {
	file_path, err := getConfigFilePath()

	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("Error converting struct to json data : %v", err)
	}

	err = os.WriteFile(file_path, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to file : %v", err)
	}

	return nil
}
