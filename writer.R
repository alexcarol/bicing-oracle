hello <- function( name ) {
    sprintf( "Hello, %s", name );
}

saveRDS(hello, "h.txt")
