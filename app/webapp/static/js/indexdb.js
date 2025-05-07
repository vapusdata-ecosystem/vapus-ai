function localDbOps(dbName, storeName, pKey) {
    const db = new LocalDb(dbName, storeName, pKey);
    return db;
}

export class LocalDb {
    constructor(dbName, storeName, pKey) {
        this.dbName = dbName;
        this.storeName = storeName;
        this.pKey = pKey;
        this.db = undefined;
        this.uniqueIndexes = [];
        this.nonUniqueIndexes = [];
    }

    // Open (or create) the database
    openIndexDB = () => {
        return new Promise((resolve, reject) => {
            const request = indexedDB.open(this.dbName, 1);
            request.onupgradeneeded = (event) => {
                this.db = event.target.result;
                // Create an object store if it doesn't already exist
                if (!this.db.objectStoreNames.contains(this.storeName)) {
                    this.db.createObjectStore(this.storeName, { keyPath: this.pKey });
                }
            };

            request.onsuccess = (event) => {
                this.db = event.target.result;
                resolve(this.db);
            };
            request.onerror = (event) => {
                console.error("Error opening database:", event.target.error);
                reject(event.target.error);
            };
        });
    };

    // Store JSON data in the database
    storeData = async (data) => {
        if (this.db === undefined) {
            this.openIndexDB();
        }
        return new Promise((resolve, reject) => {
            const dbcl = this.db;
            const transaction = dbcl.transaction(this.storeName, "readwrite");
            const store = transaction.objectStore(this.storeName);

            const request = store.add(data);
            request.onsuccess = () => resolve("Data stored successfully!");
            request.onerror = (event) => reject(event.target.error);
        });
    };

    // Update JSON data in the database
    putData = async (data) => {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction(this.storeName, "readwrite");
            const store = transaction.objectStore(this.storeName);

            const request = store.put(data);
            request.onsuccess = () => resolve("Data updated successfully!");
            request.onerror = (event) => reject(event.target.error);
        });
    };
    // Update JSON data in the database
    patchData = async (data,key) => {
        return new Promise((resolve, reject) => {
            const dbcl = this.db;
            const transaction = dbcl.transaction(this.storeName, "readwrite");
            const store = transaction.objectStore(this.storeName);
            const getRequest = store.get(key);


            getRequest.onsuccess = () => {
                const existingData = getRequest.result;

                if (existingData) {
                    // Append to the existing value array
                    const valList = JSON.parse(existingData.value);
                    valList.push(data.value[0]);
                    existingData.value = JSON.stringify(valList);
                    store.put(existingData); // Update the object in the store
                } else {
                    store.add(data);
                }

                resolve('Data stored successfully!');
            };

            getRequest.onerror = (event) => {
                reject('Error fetching existing data:', event.target.error);
            };

        });
    };

    // Retrieve JSON data by key
    retrieveData = async (key) => {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction(this.storeName, "readonly");
            const store = transaction.objectStore(this.storeName);

            const request = store.get(key);
            request.onsuccess = () => resolve(request.result);
            request.onerror = (event) => reject(event.target.error);
        });
    };

    // Retrieve all JSON data
    retrieveAllData = async () => {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction(this.storeName, "readonly");
            const store = transaction.objectStore(this.storeName);

            const request = store.getAll();
            request.onsuccess = () => resolve(request.result);
            request.onerror = (event) => reject(event.target.error);
        });
    };


    deleteRecord = async (key) => {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction(this.storeName, "readonly");
            const store = transaction.objectStore(this.storeName);

            const request = store.delete(key);
            request.onsuccess = () => resolve(request.result);
            request.onerror = (event) => reject(event.target.error);
        });
    };

    deleteRecordWithFilter = async (allowedValue) => {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction(this.storeName, 'readwrite');
            const store = transaction.objectStore(this.storeName);

            // Open a cursor to iterate through records
            const cursorRequest = store.openCursor();

            cursorRequest.onsuccess = (event) => {
                const cursor = event.target.result;
                if (cursor) {
                    // Check if the record's keyPath value matches the allowedValue
                    if (cursor.value[this.pKey] !== allowedValue) {
                        // Delete the record
                        store.delete(cursor.primaryKey);
                    }

                    // Continue to the next record
                    cursor.continue();
                } else {
                    // No more records
                    resolve();
                }
            };

            cursorRequest.onerror = (event) => {
                reject(event.target.error);
            };
        });
    }

    deleteStore = async () => {
        return new Promise((resolve, reject) => {
            const transaction = this.db.transaction(this.storeName, "readonly");
            const store = transaction.objectStore(this.storeName);

            const request = store.clear();
            request.onsuccess = () => resolve(request.result);
            request.onerror = (event) => reject(event.target.error);
        });
    };

    deleteDB = async () => {
        const deleteRequest = indexedDB.deleteDatabase(this.dbName);

        deleteRequest.onsuccess = function () {
            console.log(`Database '${this.dbName}' deleted successfully.`);
        };

        deleteRequest.onerror = function (event) {
            console.error("Error deleting database:", event.target.error);
        };

        deleteRequest.onblocked = function () {
            console.warn(
                `Database deletion for '${this.dbName}' is blocked. Close all open connections.`
            );
        };
    };
    keyExists = async (key) => {
        if (this.db === undefined) {
            this.openIndexDB();
        }
        return new Promise((resolve, reject) => {
            const dbcl = this.db;
            const transaction = dbcl.transaction(this.storeName, "readwrite");
            const store = transaction.objectStore(this.storeName);
            const getRequest = store.get(key);

            getRequest.onsuccess = (event) => {
                // If event.target.result is null, key doesn't exist
                const result = event.target.result;
                resolve(result !== undefined && result !== null);
            };

            getRequest.onerror = (event) => {
                reject(event.target.error);
            };
        });
    }

    countFilteredRecords = (lowerBound, upperBound) => {
        const request = indexedDB.open(this.dbName);

        request.onsuccess = function (event) {
            const db = event.target.result;

            // Start a transaction in "readonly" mode
            const transaction = db.transaction(this.storeName, "readonly");
            const objectStore = transaction.objectStore(this.storeName);

            // Define a key range
            const keyRange = IDBKeyRange.bound(lowerBound, upperBound);

            // Use the count() method with the key range
            const countRequest = objectStore.count(keyRange);

            countRequest.onsuccess = function () {
                console.log(
                    `Records in '${this.storeName}' between ${lowerBound} and ${upperBound}:`,
                    countRequest.result
                );
            };

            countRequest.onerror = function (event) {
                console.error("Error counting filtered records:", event.target.error);
            };
        };

        request.onerror = function (event) {
            console.error("Error opening database:", event.target.error);
        };
    };

    countRecordsByIndex = (indexName, value) => {
        const request = indexedDB.open(this.dbName);

        request.onsuccess = function (event) {
            const db = event.target.result;

            // Start a transaction in "readonly" mode
            const transaction = db.transaction(this.storeName, "readonly");
            const objectStore = transaction.objectStore(this.storeName);

            // Access the index
            const index = objectStore.index(indexName);

            // Use the count() method on the index
            const countRequest = index.count(value);

            countRequest.onsuccess = function () {
                console.log(
                    `Records in '${this.storeName}' with ${indexName} = ${value}:`,
                    countRequest.result
                );
            };

            countRequest.onerror = function (event) {
                console.error("Error counting records by index:", event.target.error);
            };
        };

        request.onerror = function (event) {
            console.error("Error opening database:", event.target.error);
        };
    };

}

// function cleanupEntries(db) {
//     return new Promise((resolve, reject) => {
//         const transaction = db.transaction("MyStore", "readwrite");
//         const store = transaction.objectStore("MyStore");

//         // Retrieve all entries
//         const request = store.getAll();
//         request.onsuccess = (event) => {
//             const allEntries = event.target.result;

//             // Sort by timestamp (latest first)
//             allEntries.sort((a, b) => b.timestamp - a.timestamp);

//             // Get entries to delete
//             const entriesToDelete = allEntries.slice(10);

//             // Delete older entries
//             for (const entry of entriesToDelete) {
//                 store.delete(entry.id);
//             }

//             resolve();
//         };
//         request.onerror = (event) => reject(event.target.error);
//     });
// }

// // Schedule cleanup every 10 minutes
// async function startCleanupRoutine() {
//     const db = await openDatabase();

//     setInterval(async () => {
//         console.log("Running cleanup...");
//         await cleanupEntries(db);
//         console.log("Cleanup completed.");
//     }, 10 * 60 * 1000); // 10 minutes
// }