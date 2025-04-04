package config

import("os")

type Config struct {
    DatabaseURL string
    
}

// LoadConfig загружает конфигурацию из переменных окружения или файла.
func LoadConfig() Config {  
    cfg := Config{
        DatabaseURL: os.Getenv("DATABASE_URL"),
        
    }
    return cfg
}