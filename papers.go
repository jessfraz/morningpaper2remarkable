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
		// The Noria Paper
		"https://www.usenix.org/conference/osdi18/presentation/gjengset": "https://pdos.csail.mit.edu/papers/noria:osdi18.pdf",
		// The FuzzyLog Paper
		"https://www.usenix.org/conference/osdi18/presentation/lockerman": "https://www.usenix.org/system/files/osdi18-lockerman.pdf",
		// Sharding the Shards peper.
		"https://www.usenix.org/conference/osdi18/presentation/annamalai": "https://www.usenix.org/system/files/osdi18-annamalai.pdf",
		// ASAP paper.
		"https://www.usenix.org/conference/osdi18/presentation/iyer": "https://www.usenix.org/system/files/osdi18-iyer.pdf",
		// Robinhood paper.
		"https://www.usenix.org/conference/osdi18/presentation/berger": "https://www.usenix.org/system/files/osdi18-berger.pdf",
		// Maelstrom paper.
		"https://www.usenix.org/conference/osdi18/presentation/veeraraghavan": "https://www.usenix.org/system/files/osdi18-veeraraghavan.pdf",
		// LegoOS paper.
		"https://www.usenix.org/conference/osdi18/presentation/shan": "https://www.usenix.org/system/files/osdi18-shan.pdf",
		// Orca paper.
		"https://www.usenix.org/conference/osdi18/presentation/bhagwan": "https://www.microsoft.com/en-us/research/uploads/prod/2018/10/Orca-OSDI.pdf",
		// REPT paper.
		"https://www.usenix.org/conference/osdi18/presentation/weidong": "https://www.usenix.org/system/files/osdi18-cui.pdf",
		// Situ paper.
		"https://www.usenix.org/conference/osdi18/presentation/huang": "https://www.cs.jhu.edu/~huang/paper/panorama-osdi18.pdf",
		// Soccer match data paper.
		"http://www.kdd.org/kdd2018/accepted-papers/view/automatic-discovery-of-tactics-in-spatio-temporal-soccer-match-data": "https://people.cs.kuleuven.be/~jesse.davis/decroos-kdd18.pdf",
		// Spacecreaft paper.
		"http://www.kdd.org/kdd2018/accepted-papers/view/detecting-spacecraft-anomalies-using-lstms-and-nonparametric-dynamic-thresh": "https://arxiv.org/pdf/1802.04431.pdf",
		// I know you'll be back paper.
		"http://www.kdd.org/kdd2018/accepted-papers/view/i-know-youll-be-back-interpretable-new-user-clustering-and-churn-prediction": "http://hanj.cs.illinois.edu/pdf/kdd18_cyang.pdf",
		// Rosetta paper.
		"http://www.kdd.org/kdd2018/accepted-papers/view/rosetta-large-scale-system-for-text-detection-and-recognition-in-images": "https://research.fb.com/wp-content/uploads/2018/10/Rosetta-Large-scale-system-for-text-detection-and-recognition-in-images.pdf",
		// Zcash paper.
		"https://www.usenix.org/conference/usenixsecurity18/presentation/kappos": "https://smeiklej.com/files/usenix18.pdf",
		// Qsym paper.
		"https://www.usenix.org/conference/usenixsecurity18/presentation/yun": "https://www.usenix.org/system/files/conference/usenixsecurity18/sec18-yun.pdf",
		// Navex paper.
		"https://www.usenix.org/conference/usenixsecurity18/presentation/alhuzali": "https://www.usenix.org/system/files/conference/usenixsecurity18/sec18-alhuzali.pdf",
		// Facebook data paper.
		"https://www.usenix.org/conference/usenixsecurity18/presentation/cabanas": "https://www.usenix.org/system/files/conference/usenixsecurity18/sec18-cabanas.pdf",
		// Cookie jar paper.
		"https://www.usenix.org/conference/usenixsecurity18/presentation/franken": "https://www.usenix.org/system/files/conference/usenixsecurity18/sec18-franken.pdf",
		// Fear the reaper paper.
		"https://www.usenix.org/conference/usenixsecurity18/presentation/scaife": "https://www.usenix.org/system/files/conference/usenixsecurity18/sec18-scaife.pdf",
		// Artistic Styles paper.
		"https://hal.inria.fr/hal-01802131v2/document": "https://arxiv.org/pdf/1805.11155.pdf",
		// Vehicle routing problems paper.
		"https://hal.inria.fr/hal-01224562/document": "https://hal.inria.fr/hal-01224562/document",
		// The data calculator paper.
		"https://stratos.seas.harvard.edu/publications/data-calculator-data-structure-design-and-cost-synthesis-first-principles-and": "https://stratos.seas.harvard.edu/files/stratos/files/datacalculator.pdf",
		// Design continums paper.
		"https://stratos.seas.harvard.edu/publications/design-continuums-and-path-toward-self-designing-key-value-stores-know-and": "https://stratos.seas.harvard.edu/files/stratos/files/selfdesign.pdf",
		// Large scale GAN paper.
		"https://openreview.net/pdf?id=B1xsqj09Fm": "https://arxiv.org/pdf/1809.11096.pdf",
		// CORALS paper.
		"https://dl.acm.org/citation.cfm?id=3290995": "https://s3-media3.fl.yelpcdn.com/assets/srv0/engineering_pages/f63a086ef2a3/assets/vendor/pdf/DSC_R09_CORALSWhoAreMyPotentialNewCustomers.pdf",
		// Slim OS paper.
		"https://www.usenix.org/conference/nsdi19/presentation/zhuo": "https://www.usenix.org/system/files/nsdi19-zhuo.pdf",
		// Datacenter topologies paper.
		"https://www.usenix.org/conference/nsdi19/presentation/zhang": "https://www.usenix.org/system/files/nsdi19-zhang.pdf",
		// Datacenter RPCs paper.
		"https://www.usenix.org/conference/nsdi19/presentation/kalia": "https://www.usenix.org/system/files/nsdi19-kalia.pdf",
		// CURP paper.
		"https://www.usenix.org/conference/nsdi19/presentation/park": "https://www.usenix.org/system/files/nsdi19-park.pdf",
		// How bad can it git?
		"https://www.ndss-symposium.org/ndss-paper/how-bad-can-it-git-characterizing-secret-leakage-in-public-github-repositories/": "https://www.ndss-symposium.org/wp-content/uploads/2019/02/ndss2019_04B-3_Meli_paper.pdf",
	}
)
