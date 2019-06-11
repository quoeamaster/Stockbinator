package crawler

type StructGenericCrawler struct {

}

// 1) read config file for rule(s) to crawl (at least url and patterns to match) based on moduleKey (stock_module-rule)
// 2) based on the key above, invoke the corresponding crawler's crawl method and scrap out the values
// 3) output the results to a repository (filestore by default or any datastorage tech e.g. elasticsearch)
func (s *StructGenericCrawler) Crawl(moduleKey string) (err error) {


	return
}
