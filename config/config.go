package config

import "os"

// Constanta for order status
const QUEUEING = "queueing"
const ACCEPTED = "accepted"
const CONFIRMED = "confirmed"
const FINISHED = "finished"
const CANCELLED = "cancelled"

// API
const API_SUBSCRIBER = ""
const API_GOPAY = "https://go-pay-sea-cfx.herokuapp.com/api/"
const API_USER = ""

var PORT = os.Getenv("PORT")

// Database
const DB_HOST = "localhost"
const DB_PORT = "5432"

const DB_USER = "joseph"
const DB_PASSWORD = "joseph"
const DB_NAME = "order"

// Environment
const ENVIRONMENT = "development"

const HEROKU = true
var DATABASE_URL = os.Getenv("DATABASE_URL")