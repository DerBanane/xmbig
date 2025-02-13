db, err = connectToDatabase(cfg.Database)
if err != nil {
	log.Fatalf("Failed to connect to database: %v", err)
	return
}
defer db.Close()
println("Connected to the Database")