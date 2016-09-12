library(RMySQL) # will load DBI as well
library(randomForest)
library(DBI)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])
from <- as.integer(args[2])
to <- as.integer(args[3])

user <- Sys.getenv("MYSQL_RAW_DATA_USER")
dbname <- Sys.getenv("MYSQL_RAW_DATA_NAME")
password <- Sys.getenv("MYSQL_RAW_DATA_PASSWORD")
host <- Sys.getenv("MYSQL_RAW_DATA_HOST")
mydb <- dbConnect(MySQL(), user=user, dbname=dbname, password=password, host=host)

query <-  sprintf(
    "SELECT bikes, slots, weather_type, UNIX_TIMESTAMP(updatetime) as updatetime FROM fit_precalculation_2 WHERE id=%d AND updatetime >= FROM_UNIXTIME(%d) AND updatetime <= FROM_UNIXTIME(%d) ORDER BY updatetime ASC",
    stationID,
    from,
    to
)
data <- dbGetQuery(mydb, query)

isCalm <- function(weather) {
    return(weather >= 3 && weather <= 8)
}
data$weather <- isCalm(data$weather)
data$dayMoment <- data$updatetime %% 86400
data$weekday <- as.POSIXlt(as.POSIXct(data$updatetime, origin="1970-01-01"))$wday

objectID <- sprintf("/tmp/station/bike/%d.fit", stationID)
fit <- readRDS(objectID)


data$predictedAmount <- predict(fit, data)
closest_bikes <- function(df, row, diff, originalFieldName) {
    return(df[which.min(abs(df$updatetime-(df$updatetime[row]-diff))), originalFieldName])
}

append_to_data <- function(df, diff, resultingFieldName, originalFieldName) {
    for (i in 1:nrow(df)) {
        df[i, resultingFieldName] <- closest_bikes(df, i, diff, originalFieldName)
    }
    return(df)
}
data <- append_to_data(data, 1800, "bikesAtMinus30", "bikes")
data <- append_to_data(data, 1800, "predictedBikesAtMinus30", "predictedAmount")
data <- append_to_data(data, 3600, "bikesAtMinus60", "bikes")
data <- append_to_data(data, 3600, "predictedBikesAtMinus60", "predictedAmount")
data <- append_to_data(data, 5400, "bikesAtMinus90", "bikes")
data <- append_to_data(data, 5400, "predictedBikesAtMinus90", "predictedAmount")
data <- append_to_data(data, 7200, "bikesAtMinus120", "bikes")
data <- append_to_data(data, 7200, "predictedBikesAtMinus120", "predictedAmount")

data$calculatedFromMinus30 <- data$bikesAtMinus30 + data$predictedAmount - data$predictedBikesAtMinus30
data$calculatedFromMinus60 <- data$bikesAtMinus60 + data$predictedAmount - data$predictedBikesAtMinus60
data$calculatedFromMinus90 <- data$bikesAtMinus90 + data$predictedAmount - data$predictedBikesAtMinus90
data$calculatedFromMinus120 <- data$bikesAtMinus120 + data$predictedAmount - data$predictedBikesAtMinus120

write.csv(data, sprintf("station%d.csv", stationID))

dbDisconnect(mydb)
