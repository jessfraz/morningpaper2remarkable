package main

var (
	// TODO: this is horrible find a better way of downloading acm papers.
	knownPapersDownloadLinks = map[string]string{
		// The FairSwap Paper.
		"https://dl.acm.org/citation.cfm?id=3243857": "https://eprint.iacr.org/2018/740.pdf",
		// The RapidChain Paper.
		"https://dl.acm.org/citation.cfm?id=3243853": "https://eprint.iacr.org/2018/460.pdf",
		// Log Analysis paper.
		"https://dl.acm.org/citation.cfm?id=3236083": "https://microsoft.com/en-us/research/uploads/prod/2018/06/Identifying-Impactful-Service-System-Problems-via-Log-Analysis.pdf",
		// Facebook paper.
		"https://research.fb.com/publications/applied-machine-learning-at-facebook-a-datacenter-infrastructure-perspective/": "https://research.fb.com/wp-content/uploads/2017/12/hpca-2018-facebook.pdf",
	}
)
