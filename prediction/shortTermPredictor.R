library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])
predictBikes <- as.logical(args[2])
predictionTime <- as.integer(args[3])
currentTime <- as.integer(args[6])
predictedWeather <- as.integer(args[4])
currentWeather <- as.integer(args[5])

isCalm <- function(weather) {
    return(weather >= 3 && weather <= 8)
}

objectID <- sprintf("/tmp/station/bike/%d.fit", stationID)
fit <- readRDS(objectID)

currentWeekday <- as.POSIXlt(as.POSIXct(currentTime, origin="1970-01-01"))$wday
predictionWeekday <- as.POSIXlt(as.POSIXct(predictionTime, origin="1970-01-01"))$wday

weather_type <- c(isCalm(predictedWeather), isCalm(currentWeather))
weekday <- c(predictionWeekday, currentWeekday)
dayMoment <- c(predictionTime %% 86400, currentTime %% 86400)
updatetime <- c(predictionTime, currentTime)
object <- data.frame(updatetime, dayMoment, weekday, weather_type)


predict(fit, object)
