package config

import "os"

// Constanta for order status
const FINISHED = "finished"
const QUEUEING = "queueing"
const ACCEPTED = "accepted"
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
const ENVIRONMENT = "test"

var DATABASE_URL = os.Getenv("DATABASE_URL")
var REDIS_URL = "[redis://h:pebfe2bd00d785fc79f982c9aaa72d063daf874e8c864e14ef52ef7e682213091@ec2-18-214-19-12.compute-1.amazonaws.com]:26919"