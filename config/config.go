package config

import "os"

// Constanta for order status
const QUEUEING = "queueing"
const ACCEPTED = "accepted"
const CONFIRMED = "confirmed"
const FINISHED = "finished"
const CANCELLED = "cancelled"

// API
const API_SUBSCRIBER = "https://goridepay-driverworker.herokuapp.com"
const API_GOPAY = "https://go-pay-sea-cfx.herokuapp.com"

var PORT = os.Getenv("PORT")

// Database
const DB_HOST = "localhost"
const DB_PORT = "5432"

const DB_USER = "joseph"
const DB_PASSWORD = "joseph"
const DB_NAME = "order"

// Environment
const ENVIRONMENT = "development"

var HEROKU = os.Getenv("HEROKU")
var DATABASE_URL = os.Getenv("DATABASE_URL")