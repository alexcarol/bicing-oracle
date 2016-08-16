library(RMySQL) # will load DBI as well
library(randomForest)

args <- commandArgs(trailingOnly = TRUE)
stationID <-as.integer(args[1])
print("stationID")
print(stationID)

mydb <- dbConnect(MySQL(), user='root', dbname='bicing_oracle_raw')

objectID <- sprintf("station_%d.fit", stationID)

##TODO use a normalized table/view with the weather as well as the pbikes
query <-  sprintf("SELECT object FROM fits WHERE id='%s' LIMIT 1", objectID)
data <- dbGetQuery(mydb, query)

print("data[0]")
object <- data[1,]

unserializedObject <- unserialize(object)
print(unserializedObject)

dbDisconnect(mydb)
