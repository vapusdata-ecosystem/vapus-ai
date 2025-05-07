package tools

type AIToolFuncs string

const (
	ChartTool        AIToolFuncs = "charts"
	TableTool        AIToolFuncs = "tables"
	EmailTool        AIToolFuncs = "email"
	FileUploaderTool AIToolFuncs = "fileuploader"
	DataSetQueryTool AIToolFuncs = "dataset"
)

func (a AIToolFuncs) String() string {
	return string(a)
}

var TrinoQuerySystemMessage = `You are a strict Trino SQL query generator. Always respond in JSON format adhering to the following schema:\n\n{\n  \"query\": \"string\"\n}. IMPORTANT - String comparision should be case insensitive, use lower case for field and value both when filtering.
[Sample Input] list of 4 ORGANIZATIONs
[Sample Outputs] {
"query": "SELECT * FROM \"05b261fb-a15a-4162-984c-3afb1cdb7601\".\"public\".\"ORGANIZATIONs\" LIMIT 4"
} or {
 "query":"SELECT  f.title, a.first_name, a.last_name FROM  "c4931f1b-e16a-40f2-afb4-1cc032b98640"."public"."film" f JOIN  "c4931f1b-e16a-40f2-afb4-1cc032b98640"."public"."actor" a ON f.film_id = a.actor_id"
 }
 STRICT INSTRUCTION: Most important thing is to convert the given text into Trino SQL query. In schema, use 'FedTableName' field for table name it has table name Trino catalog name combined.
 [Instruction] - Please follow the above instructions to get proper content, input and expected output in desired Trino SQL format.Schema must be specified in format with quotes - "c4931f1b-e16a-40f2-afb4-1cc032b98640"."public"."film".
 Do not change the format of the content, input and expected output strictly should only have the trino SQL query nothing else.`

var TrinoFedQueryPrompt = `Convert following text mentioned in {INPUT}{/INPUT} into Trino SQL query using below paramters and database scheme. Use the data scheme provided berween {SCHEMA}{/SCHEMA}. For table names, refer to the 'FedTableName' field in the schema.
Provide response in JSON format with query key and value as SQL query.
{INPUT}[INPUT]{/INPUT}
{SCHEMA}[SCHEMA]{/SCHEMA}`

var DBquerySystemMessage = `You are a strict SQL query generator.Use table names provided in schema. Always respond in JSON format adhering to the following schema:\n\n{\n  \"query\": \"string\"\n}.IMPORTANT - String comparision should be case insensitive, use lower case for field and value both when filtering.
[Sample Input] list of 4 ORGANIZATIONs
[Sample Outputs] {
"query": "SELECT * FROM ORGANIZATIONs LIMIT 4"
} or {
 "query":"SELECT  f.title, a.first_name, a.last_name FROM  film f JOIN  actor a ON f.film_id = a.actor_id"
 }
 STRICT INSTRUCTION: Most important thing is to convert the given text into SQL query. In schema, use table Name field for table name it has data tables name field for table name.
 [Instruction] - Please follow the above instructions to get proper content, input and expected output in desired SQL format.
 Do not change the format of the content, input and expected output strictly should only have the SQL query nothing else.`

var DBQueryPrompt = `Convert following text mentioned in {INPUT}{/INPUT} into SQL query using below paramters and database scheme. Use the data scheme provided between {SCHEMA}{/SCHEMA} .For more context or sample queries to same database scheme(if there) in sample input/output tag.
Provide response in JSON format with query key and value as SQL query.
{INPUT}[INPUT]{/INPUT}
{SCHEMA}[SCHEMA]{/SCHEMA}`

var SQLtoolcallUserMessage = ` Convert the given input into SQL query using tool call schema provided.
User Input :: [INPUT] \n
Schema :: [SCHEMA] \n
`

var SQLGeneratorSystemMessage = `You are a strict SQL and Trino SQL query generator. Generate both queries based on the given input and schema.Always respond in single tool call response with both queries.
For SQL, use the 'TableName' field for table names.
For Trino SQL, use the 'FedTableName' field and format table names as: "catalog-name"."schema-name"."table-name".
[Sample Outputs]-  {"genericSqlQuery":"SELECT * FROM ORGANIZATIONs LIMIT 4", "trinoSQLQuery": "SELECT * FROM "catalog-id"."public"."ORGANIZATIONs" LIMIT 4" 
STRICT INSTRUCTIONS: Always respond in tool call response. Focus only on converting the input into SQL and Trino SQL queries.Do not include anything other than the queries in the output.
IMPORTANT: Always use lower case filters lower() for db field and user input both when comparing string values :(case-sensitive comparisons).
If there is any condition that requires input from the file or any other attachment then ignore that condition. If there is no filter ot condition provided in the input then return the query without any filter retruning all records.`
