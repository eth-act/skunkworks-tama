package zkvm

// This file is for the host. It has been intentionally left blank so I can write notes in here.
/*
	When the host sends data to the guest, the general idea is that they will serialize that data using serialization.MustSerialize
	Then they will send it to Zisk by writing that data into an inputs file.
*/
