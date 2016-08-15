library(RMySQL) # will load DBI as well
library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])

mydb <- dbConnect(MySQL(), user='root', dbname='bicing_oracle_raw')

##TODO use a normalized table/view with the weather as well as the pbikes
query <-  sprintf(
    "SELECT pbikes, pslots, UNIX_TIMESTAMP(updatetime) as updatetime FROM fit_precalculation WHERE id=%d", ## updatetime > X?
    stationID
)
data <- dbGetQuery(mydb, query)

fit <- randomForest(as.factor(pbikes) ~ updatetime, data=data, importance=TRUE, ntree=100)

dbGetQuery(mydb, "CREATE TABLE IF NOT EXISTS fits (id varchar(255), object text)")

serializedObject <- serialize(fit, NULL, ascii=T)

objectID <- sprintf("station_%d.fit", stationID)
dbGetQuery(mydb, sprintf("INSERT INTO fits VALUES('%s', '%s')", objectID, fit))

##TODO use a normalized table/view with the weather as well
dbDisconnect(mydb)
