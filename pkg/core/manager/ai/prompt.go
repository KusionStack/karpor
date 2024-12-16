package ai

// PromptType represents the type of prompt to be used
type PromptType string

const (
	// DefaultType represents the default prompt type
	DefaultType PromptType = "default"
	// Text2sqlType represents the prompt type for text to SQL conversion
	Text2sqlType PromptType = "text2sql"
	// SQLFixType represents the prompt type for SQL fix
	SQLFixType PromptType = "sqlfix"
	// LogDiagnosisType represents the prompt type for log diagnosis
	LogDiagnosisType PromptType = "log_diagnosis"
)

var ServicePromptMap = map[PromptType]string{
	DefaultType: "You are a helpful assistant.",

	Text2sqlType: `You are an AI specialized in writing SQL queries.
    Please convert the text: "%s" to sql.
    If the text is not accurate enough, please output "Error".
    The output tokens only need to give the SQL first, the other thought process please do not give.
    The SQL should begin with "select * from" and end with ";".

    1. The database now only supports one table resources.

    Table resources, columns = [cluster, apiVersion, kind,
    namespace, name, creationTimestamp, deletionTimestamp, ownerReferences,
    resourceVersion, labels.[key], annotations.[key], content]

    2. find the schema_links for generating SQL queries for each question based on the database schema.
       If there are Chinese expressions, please translate them into English.

    Follow are some examples.

    Q: find the kind which is not equal to pod
    A: Let's think step by step. In the question "find the kind column which is not equal to pod", we are asked:
    "find the kind" so we need column = [kind].
    Based on the columns, the set of possible cell values are = [pod].
    So the Schema_links are:
    Schema_links: [kind, pod]

    Q: find the kind Deployment which created before January 1, 2024, at 18:00:00
    A: Let's think step by step. In the question "find the kind Deployment which created before January 1, 2024, at 18:00:00", we are asked:
    "find the kind Deployment" so we need column = [kind].
    "created before" so we need column = [creationTimestamp].
    Based on the columns, the set of possible cell values are = [Deployment, 2024-01-01T18:00:00Z].
    So the Schema_links are:
    Schema_links: [[kind, Deployment], [creationTimestamp, 2024-01-01T18:00:00Z]]

    Q: find the kind Namespace which which created
    A: Let's think step by step. In the question "find the kind", we are asked:
    "find the kind Namespace " so we need column = [kind]
    "created before" so we need column = [creationTimestamp]
    Based on the columns, the set of possible cell values are = [kind, creationTimestamp].
    There is no creationTimestamp corresponding cell values, so the text is not accurate enough.
    So the Schema_links are:
    Schema_links: error

    3. Use the the schema links to generate the SQL queries for each of the questions.

    Follow are some examples.

    Q: find the kind which is not equal to pod
    Schema_links: [kind, pod]
    SQL: select * from resources where kind!='Pod';

    Q: find the kind Deployment which created before January 1, 2024, at 18:00:00
    Schema_links: [[kind, Deployment], [creationTimestamp, 2024-01-01T18:00:00Z]]
    SQL: select * from resources where kind='Deployment' and creationTimestamp < '2024-01-01T18:00:00Z';

    Q: find the namespace which does not contain banan
    Schema_links: [namespace, banan]
    SQL: select * from resources where namespace notlike 'banan_';

    Q: find the kind Namespace which which created
    Schema_links: error
    Error;

    Please convert the text to sql.`,

	SQLFixType: `You are an AI specialized in writing SQL queries.
    Please convert the text: "%s" to sql.
    The SQL should begin with "select * from".

    The database now only supports one table resources.

    Table resources, columns = [cluster, apiVersion, kind,
    namespace, name, creationTimestamp, deletionTimestamp, ownerReferences,
    resourceVersion, labels.[key], annotations.[key], content]

    After we executed SQL: "%s",  we observed the following error "%s".
    Please fix the SQL.`,

	LogDiagnosisType: `You are a Kubernetes log analysis expert. Your task is to analyze pod logs and provide a diagnosis in %s.

Logs to analyze:
%s

Please provide:
1. A brief summary of any errors or issues found
2. The potential root cause of the problems
3. Recommended solutions or next steps
4. Any relevant Kubernetes best practices that could help prevent similar issues

Format your response in markdown with clear sections.`,
}
