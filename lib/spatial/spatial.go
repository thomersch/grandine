package spatial

type PropertyRetriever interface {
	Properties() map[string]interface{}
}
