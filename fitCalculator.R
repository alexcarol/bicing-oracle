install.packages("RMySQL")

calculateFit <- function(stationID){
    library(RMySQL) # will load DBI as well
    mydb <- dbConnect(MySQL(), user='user', password='password', dbname='bicing_raw', host='mysql_raw')
    ##con <- dbConnect(dbDriver("MySQL"), dbname = "test")

    // TODO use a normalized table/view with the weather as well
    query <-  sprintf(
        "SELECT bikes, slots, UNIX_TIMESTAMP(updatetime) FROM station_state WHERE id=%d",
        stationID
    )
    data <- dbGetQuery(con, query)

    library(randomForest)
    fit <- randomForest(as.factor(pbikes) ~ Day.moment + Lunes + Martes + Mi.rcoles + Jueves + Viernes + S.bado + time,
                    data=data,
                    importance=TRUE, 
                    ntree=100)

    dbDisconnect(con)
}
