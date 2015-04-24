# sdb
#### Simple database for data blocks(i.e. time series) with key
It provides write port (persist with key and data w/o timestamp) and read port (query with key w/o timestamp).
* Write port
	* append or delete data with key on change log with index update
	* upon reaching threshold merge change log onto main storage file and update index
	* append or delete data don't lead to data deletion directly. Only merge operation squeezes out obsolete data
* Read port
	* secondary index (one level B*-tree) in memory
	* query with key by search secondary index and index for the offset and length of the corresponding data block
	* optionally filter data subblocks using timestamp
	* monitor working sets (main storage file, change logs) to synchronize with storage

Notice	
* It is meant to be used to store data that are not updated frequently, e.g. for hourly, dayly etc.. rolling
* Queries are not parallel or concurrent, which prevents intervening seeking operations that lead to poor performance
