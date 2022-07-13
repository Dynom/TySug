package service

import (
	"context"
	"testing"

	"github.com/Dynom/TySug/finder"

	"github.com/sirupsen/logrus/hooks/test"
)

func TestNewDomainWithError(t *testing.T) {
	log, _ := test.NewNullLogger()

	_, err := NewDomain([]string{}, log, finder.WithAlgorithm(nil))

	if err == nil {
		t.Error("Expecting an error to have been thrown.")
	}
}

func TestServiceFindSimilarWords(t *testing.T) {
	l, _ := test.NewNullLogger()

	testData := []struct {
		Expect     string
		Inputs     []string
		References []string
	}{
		{
			// On Query-US, Left hand fingers are more likely to travel north/north-east on these inputs
			Expect: "beer",
			Inputs: []string{"bee4", "bee5"},
			References: []string{
				"beefinesses", "beekeepings", "beerinesses", "beechdrops", "beechmasts", "beechwoods", "beefeaters",
				"beefsteaks", "beekeepers", "beekeeping", "beebreads", "beechiest", "beechmast", "beechnuts", "beechwood",
				"beefaloes", "beefcakes", "beefeater", "beefiness", "beefsteak", "beefwoods", "beekeeper", "beelining",
				"beeriness", "beestings", "beeswaxes", "beeswings", "beetroots", "beebread", "beechier", "beechnut", "beefalos",
				"beefcake", "beefiest", "beefless", "beefwood", "beehives", "beelined", "beelines", "beeriest", "beeswing",
				"beetlers", "beetling", "beetroot", "beeyards", "beebees", "beechen", "beeches", "beedies", "beefalo",
				"beefier", "beefily", "beefin", "beehiv", "beelik", "beelin", "beepers", "beeping", "beerier", "beeswax",
				"beetled", "beetler", "beetles", "beeyard", "beezers", "beebee", "beechy", "beefed", "beeped", "beeper",
				"beetle", "beeves", "beezer", "beech", "beedi", "beefs", "beefy", "beeps", "beers", "beery", "beets", "beef",
				"been", "beep", "beer", "bees", "beet",

				// Jaro Winkler doesn't work well on -1 penalties
				//"bee",
			},
		},
		{
			Expect: "beef",
			Inputs: []string{
				"beed",
			},
			References: []string{
				"beef", "been", "beep", "beer", "bees", "beet",
			},
		},
	}

	for _, td := range testData {

		svc, _ := NewDomain(td.References, l)
		for _, input := range td.Inputs {
			result, score, _ := svc.Find(context.Background(), input)

			if result != td.Expect {
				t.Errorf("Expected the result to be %s, instead I got %s.", td.Expect, result)
				t.Logf("\nInput : %s\nResult: %s\nScore : %f", input, result, score)
			}
		}
	}
}

func TestCommonTypos(t *testing.T) {
	list := []string{}

	l, _ := test.NewNullLogger()
	svc, _ := NewDomain(list, l)

	testData := []struct {
		Input  string
		Expect string
	}{
		{},
	}

	for _, td := range testData {
		result, score, _ := svc.Find(context.Background(), td.Input)

		if td.Expect != result {
			t.Errorf("Expected input %s to result in %s, instead I got %s (score: %f)",
				td.Input, td.Expect, result, score,
			)
		}
	}
}

func TestExpectations(t *testing.T) {
	l, _ := test.NewNullLogger()

	testData := []struct {
		Inputs []string
		Expect string
	}{
		// These inputs do not match well against the list of domains. This is largely because of the many different
		// tld's of the Google domains.
		//{Expect: "google.com", Inputs: []string{"oogle.com", "gogle.com", "goole.com", "google.om"}},
		{Expect: "google.com", Inputs: []string{
			"foogle.com", "voogle.com", "boogle.com", "hoogle.com", "yoogle.com", "toogle.com", "giogle.com", "gkogle.com",
			"glogle.com", "gpogle.com", "g0ogle.com", "g9ogle.com", "goigle.com", "gokgle.com", "golgle.com", "gopgle.com",
			"go0gle.com", "go9gle.com", "goofle.com", "goovle.com", "gooble.com", "goohle.com", "gooyle.com", "gootle.com",
			"googke.com", "googpe.com", "googoe.com", "googlw.com", "googls.com", "googld.com", "googlr.com", "googl4.com",
			"googl3.com", "google.xom", "google.vom", "google.fom", "google.dom", "google.cim", "google.ckm",
			// "google.clm", // fails, matches with 'google.cl'. @todo should shorter variants be penalised harder?
			"google.cpm", "google.c0m", "google.c9m", "google.con", "google.cok", "google.coj", "googe.com", "googl.com",
			"googlecom", "google.cm", "google.co", "ogogle.com", "google.com", "gogole.com", "goolge.com", "googel.com",
			"googl.ecom", "googlec.om", "google.ocm", "google.cmo", "ggoogle.com", "gooogle.com", "gooogle.com",
			"googgle.com", "googlle.com", "googlee.com", "google..com", "google.ccom", "google.coom", "google.comm",
		}},

		// These inputs do not match well against the list of domains.
		//{Expect: "netflix.com", Inputs: []string{"etflix.com", "ntflix.com"}},
		{Expect: "netflix.com", Inputs: []string{
			"betflix.com", "metflix.com", "jetflix.com", "hetflix.com", "nwtflix.com", "nstflix.com", "ndtflix.com",
			"nrtflix.com", "n4tflix.com", "n3tflix.com", "nerflix.com", "nefflix.com", "negflix.com", "neyflix.com",
			"ne6flix.com", "ne5flix.com", "netdlix.com", "netclix.com", "netvlix.com", "netglix.com", "nettlix.com",
			"netrlix.com", "netfkix.com", "netfpix.com", "netfoix.com", "netflux.com", "netfljx.com", "netflkx.com",
			"netflox.com", "netfl9x.com", "netfl8x.com", "netfliz.com", "netflic.com", "netflid.com", "netflis.com",
			"netflix.xom", "netflix.vom", "netflix.fom", "netflix.dom", "netflix.cim", "netflix.ckm", "netflix.clm",
			"netflix.cpm", "netflix.c0m", "netflix.c9m", "netflix.con", "netflix.cok", "netflix.coj",
			"neflix.com", "netlix.com", "netfix.com", "netflx.com", "netfli.com", "netflixcom", "netflix.om",
			"netflix.cm", "netflix.co", "entflix.com", "nteflix.com", "neftlix.com", "netlfix.com", "netfilx.com",
			"netflxi.com", "netfli.xcom", "netflixc.om", "netflix.ocm", "netflix.cmo", "nnetflix.com", "neetflix.com",
			"nettflix.com", "netfflix.com", "netfllix.com", "netfliix.com", "netflixx.com", "netflix..com", "netflix.ccom",
			"netflix.coom", "netflix.comm",
		}},

		{Expect: "4399.com", Inputs: []string{
			"3399.com", "e399.com", "r399.com", "5399.com", "4299.com", "4w99.com", "4e99.com", "4499.com", "4389.com",
			"43i9.com", "43o9.com", "4309.com", "4398.com", "439i.com", "439o.com", "4390.com", "4399.xom", "4399.vom",
			"4399.fom", "4399.dom", "4399.cim", "4399.ckm", "4399.clm", "4399.cpm", "4399.c0m", "4399.c9m", "4399.con",
			"4399.cok", "4399.coj", "399.com", "499.com", "439.com", "439.com", "4399com", "4399.om", "4399.cm", "4399.co",
			"3499.com", "4939.com", "4399.com", "439.9com", "4399c.om", "4399.ocm", "4399.cmo", "44399.com", "43399.com",
			"43999.com", "43999.com", "4399..com", "4399.ccom", "4399.coom", "4399.comm",
		}},

		{Expect: "yahoo.com", Inputs: []string{
			"yaoo.com", "yajoo.com",
		}},
	}

	svc, _ := NewDomain(domains, l)
	for _, td := range testData {
		for _, input := range td.Inputs {
			result, score, _ := svc.Find(context.Background(), input)

			if td.Expect != result {
				t.Errorf("Expected input %s to result in %s, instead I got %s (score: %f)",
					input, td.Expect, result, score,
				)
			}
		}
	}
}

var domains = []string{
	"google.com", "youtube.com", "facebook.com", "baidu.com", "wikipedia.org", "yahoo.com", "qq.com", "amazon.com",
	"taobao.com", "twitter.com", "tmall.com", "instagram.com", "google.co.in", "live.com", "sohu.com", "vk.com", "jd.com",
	"reddit.com", "sina.com.cn", "weibo.com", "google.co.jp", "360.cn", "login.tmall.com", "blogspot.com", "linkedin.com",
	"yandex.ru", "google.ru", "google.co.uk", "netflix.com", "google.com.hk", "google.com.br", "yahoo.co.jp",
	"csdn.net", "t.co", "microsoft.com", "ebay.com", "google.fr", "google.de", "alipay.com", "pages.tmall.com", "twitch.tv",
	"msn.com", "office.com", "xvideos.com", "bing.com", "mail.ru", "aliexpress.com", "stackoverflow.com", "naver.com",
	"github.com", "livejasmin.com", "whatsapp.com", "imgur.com", "wikia.com", "microsoftonline.com", "google.it",
	"google.ca", "amazon.co.jp", "tumblr.com", "google.com.tw", "google.com.mx", "imdb.com", "xhamster.com",
	"wordpress.com", "google.es", "google.com.tr", "paypal.com", "adobe.com", "tribunnews.com", "thestartmagazine.com",
	"google.co.kr", "apple.com", "google.com.au", "bilibili.com", "pinterest.com", "googleusercontent.com", "xnxx.com",
	"hao123.com", "bongacams.com", "txxx.com", "booking.com", "fbcdn.net", "bbc.co.uk", "cobalten.com", "coccoc.com",
	"quora.com", "bbc.com", "soso.com", "amazon.in", "amazon.de", "pixnet.net", "xinhuanet.com", "detail.tmall.com",
	"google.co.id", "gmw.cn", "dropbox.com", "amazon.co.uk", "espn.com", "google.pl", "craigslist.org", "mozilla.org",
	"google.co.th", "detik.com", "tianya.cn", "zhihu.com", "vimeo.com", "dailymail.co.uk", "so.com", "chaturbate.com",
	"google.com.ua", "google.com.ar", "google.com.pk", "soundcloud.com", "cnn.com", "google.com.sa", "amazonaws.com",
	"tokopedia.com", "ask.com", "google.com.eg", "theguardian.com", "onlinesbi.com", "force.com", "stackexchange.com",
	"1688.com", "rakuten.co.jp", "nytimes.com", "iqiyi.com", "chase.com", "youm7.com", "roblox.com", "aparat.com",
	"dailymotion.com", "popads.net", "clearload.bid", "speakol.com", "spotify.com", "indeed.com", "nicovideo.jp",
	"sogou.com", "discordapp.com", "salesforce.com", "avito.ru", "nih.gov", "panda.tv", "bukalapak.com",
	"google.co.za", "daum.net", "google.com.vn", "ok.ru", "dkn.tv", "globo.com", "exosrv.com", "openload.co",
	"w3schools.com", "google.az", "google.co.ao", "mediafire.com", "thepiratebay.org", "ettoday.net", "ebay.de",
	"cnet.com", "softonic.com", "deviantart.com", "ebay.co.uk", "cnblogs.com", "google.co.ve", "yelp.com", "fc2.com",
	"godaddy.com", "liputan6.com", "alibaba.com", "uol.com.br", "vice.com", "mercadolivre.com.br", "zhanqi.tv",
	"slideshare.net", "google.cn", "google.gr", "etsy.com", "walmart.com", "google.nl", "google.com.sg", "flipkart.com",
	"huanqiu.com", "onlinevideoconverter.com", "wikihow.com", "zillow.com", "google.com.co", "google.com.ph", "hulu.com",
	"amazon.it", "wetransfer.com", "fifa.com", "douban.com", "albawabhnews.com", "downloadhelper.net", "amazon.fr",
	"mama.cn", "buzzfeed.com", "google.com.pe", "tripadvisor.com", "china.com.cn", "neveryone.club", "163.com", "setn.com",
	"caijing.com.cn", "nownews.com", "rambler.ru", "aliyun.com", "foxnews.com",
	"indiatimes.com", "spankbang.com", "bankofamerica.com", "kompas.com", "xfinity.com", "shutterstock.com",
	"sciencedirect.com", "slack.com", "forbes.com", "foxsports.com", "wellsfargo.com", "crptentry.com", "ebc.net.tw",
	"mi.com", "metropcs.mobi", "google.be", "eastday.com", "speedtest.net", "twimg.com", "blogger.com", "digikala.com",
	"doubleclick.net", "trello.com", "washingtonpost.com", "sindonews.com", "google.pt", "ikea.com", "scribd.com",
	"researchgate.net", "savefrom.net", "yts.am", "freepik.com", "steamcommunity.com", "livejournal.com", "google.com.ng",
	"google.cl", "mega.nz", "exerciers.mobi", "google.se", "obozrevatel.com", "babytree.com", "ladbible.com", "youth.cn",
	"zendesk.com", "sberbank.ru", "amazon.es", "ebay-kleinanzeigen.de", "canva.com", "myshopify.com", "medium.com",
	"genius.com", "google.ro", "okezone.com", "ci123.com", "gamespot.com", "bet9ja.com", "diply.com", "amazon.ca",
	"outbrain.com", "duckduckgo.com", "myway.com", "ltn.com.tw", "hclips.com", "messenger.com", "momoshop.com.tw",
	"weather.com", "getmyapp1.com", "gamepedia.com", "rednet.cn", "hdzog.com", "breitbart.com", "banvenez.com",
	"haber7.com", "youdao.com", "amazon.cn", "rt.com", "youku.com", "uzone.id", "google.ae", "jf71qh5v14.com",
	"archive.org", "nextoptim.com", "baylnk.com", "varzesh3.com", "allegro.pl", "sabah.com.tr", "oath.com", "usps.com",
	"google.ch", "steampowered.com", "siteadvisor.com", "airbnb.com", "kakaku.com", "cloudfront.net", "patria.org.ve",
	"google.co.il", "abs-cbn.com", "hibids10.com", "ups.com", "gfycat.com", "tistory.com", "google.dz",
	"businessinsider.com", "exdynsrv.com", "kinopoisk.ru", "1337x.to", "mailchimp.com", "fedex.com", "ign.com",
	"huffingtonpost.com", "intuit.com", "google.at", "glassdoor.com", "zoho.com", "aol.com", "toutiao.com", "google.cz",
	"yy.com", "tvbs.com.tw", "telegram.org", "bestbuy.com", "marca.com", "goodreads.com", "rutracker.org", "gismeteo.ru",
	"hp.com", "ameblo.jp", "line.me", "homedepot.com", "ensonhaber.com", "google.kz", "behance.net", "uptodown.com",
	"atlassian.net", "kaskus.co.id", "bles.com", "jianshu.com", "bet365.com", "slickdeals.net", "zoom.us", "box.com",
	"office365.com", "pixiv.net", "patreon.com", "zippyshare.com", "udemy.com", "capitalone.com", "hotstar.com",
	"hotmovs.com", "americanexpress.com", "gearbest.com", "hola.com", "wikimedia.org", "onet.pl",
	"incometaxindiaefiling.gov.in", "wp.pl", "reverso.net", "oracle.com", "friv.com", "coinmarketcap.com", "asos.com",
	"dell.com", "sourceforge.net", "livedoor.jp", "leboncoin.fr", "thesaurus.com", "icloud.com", "olx.ua", "ngacn.cc",
	"battle.net", "list.tmall.com", "irctc.co.in", "mercadolibre.com.ar", "4chan.org", "flickr.com", "myanimelist.net",
	"okta.com", "webex.com", "ifeng.com", "kissanime.ru", "hdfcbank.com", "ndtv.com", "51sole.com", "dmm.co.jp",
	"wordpress.org", "exoclick.com", "ouo.io", "skype.com", "independent.co.uk", "bloomberg.com", "google.ie",
	"instructure.com", "samsung.com", "googlevideo.com", "newsprofin.com", "mercadolibre.com.mx", "accuweather.com",
	"wordreference.com", "telewebion.com", "google.hu", "usatoday.com", "9gag.com", "fiverr.com", "target.com",
	"jrj.com.cn", "mlb.com", "google.no", "orange.fr", "redd.it", "pinimg.com", "telegraph.co.uk", "elpais.com",
	"android.com", "yao.tmall.com", "cdninstagram.com", "17ok.com", "wix.com", "epicgames.com", "rarbg.to", "indoxx1.com",
	"kapanlagi.com", "okdiario.com", "taleo.net", "wiley.com", "manoramaonline.com", "vtv.vn", "taboola.com",
	"cambridge.org", "giphy.com", "1bcde.com", "debate.com.mx", "pixabay.com", "kooora.com", "weebly.com", "people.com.cn",
	"patch.com", "doublepimpssl.com", "primevideo.com", "op.gg", "goal.com", "uidai.gov.in", "beeg.com", "spiegel.de",
	"evernote.com", "google.sk", "shopify.com", "jb51.net", "shaparak.ir", "ria.ru", "upwork.com", "themeforest.net",
	"bp.blogspot.com", "go.com", "rediff.com", "livedoor.com", "gmx.net", "userapi.com", "seznam.cz", "free.fr",
	"gizmodo.com", "theverge.com", "icicibank.com", "vidio.com", "taringa.net", "playstation.com", "springer.com",
	"lenta.ru", "cnbc.com", "urbandictionary.com", "goo.ne.jp", "perfecttoolmedia.com", "adp.com", "spotscenered.info",
	"express.co.uk", "as.com", "rottentomatoes.com", "sahibinden.com", "investing.com", "tahiamasr.com", "olx.pl",
	"news.com.au", "google.fi", "google.dk", "files.wordpress.com", "elsevier.com", "merdeka.com", "list-manage.com",
	"huaban.com", "pikabu.ru", "rotumal.com", "billdesk.com", "wixsite.com", "nike.com", "wish.com", "dcinside.com",
	"shein.com", "google.com.kw", "nextlnk1.com", "asus.com", "ebay.com.au", "google.co.nz", "citi.com", "utorrent.com",
	"expedia.com", "ebay.it", "subscene.com", "zdf.de", "kickstarter.com", "grid.id", "wowhead.com",
	"gsmarena.com", "weblio.jp", "3dmgame.com", "wiktionary.org", "divar.ir", "hm.com", "104.com.tw", "libero.it",
	"hurriyet.com.tr", "engadget.com", "memurlar.net", "doublepimp.com", "nexusmods.com", "emol.com", "sonyliv.com",
	"repubblica.it", "58.com", "gamersky.com", "ouedkniss.com", "tutorialspoint.com", "ibm.com", "xda-developers.com",
	"reuters.com", "surveymonkey.com", "blackboard.com", "att.com", "mit.edu", "drom.ru", "kumparan.com", "mmoframes.com",
	"naukri.com", "xiaomi.com", "hellomagazine.com", "naija.ng", "paytm.com", "informationvine.com", "discogs.com",
	"wsj.com", "smadav.net", "hespress.com", "drudgereport.com", "yenisafak.com", "cbssports.com", "bitly.com", "ck101.com",
	"stockstar.com", "google.lk", "cisco.com", "souq.com", "nur.kz", "filehippo.com", "idntimes.com", "bandcamp.com",
	"bitauto.com", "appledaily.com", "hupu.com", "web.de", "quizlet.com", "groupon.com", "elbalad.news", "tube8.com",
	"banggood.com", "biobiochile.cl", "blog.jp", "pantip.com", "autodesk.com", "nhk.or.jp", "smallpdf.com", "google.by",
	"prothomalo.com", "sabq.org", "douyu.com", "inquirer.net", "ultimate-guitar.com", "gamer.com.tw", "convert2mp3.net",
	"gyazo.com", "eyny.com", "runoob.com", "hotels.com", "wattpad.com", "eventbrite.com", "umblr.com", "tradingview.com",
	"cnnindonesia.com", "bleacherreport.com", "healthline.com", "nypost.com", "qingdaonews.com", "gmarket.co.kr",
	"nyaa.si", "ieee.org", "ca.gov", "bancodevenezuela.com", "fidelity.com", "gosuslugi.ru", "chinaz.com",
	"squarespace.com", "oload.download", "grammarly.com", "webmd.com", "cinecalidad.to", "viralvideos.technology",
	"namnak.com", "blatungo.com", "google.com.do", "bild.de", "crunchyroll.com", "rbc.ru", "corriere.it", "drive2.ru",
	"abola.pt", "udn.com", "browser-tools.systems", "discuss.com.hk", "coursera.org", "foxsportsgo.com", "smzdm.com",
	"rarbg.is", "milliyet.com.tr", "alicdn.com", "ptt.cc", "thesun.co.uk", "championat.com", "lifehacker.com",
	"epochtimes.com", "blog.me", "pandora.com", "thefreedictionary.com", "jw.org", "zhaopin.com",
	"vnexpress.net", "zhibo8.cc", "gstatic.com", "banesconline.com", "merriam-webster.com", "interia.pl", "techradar.com",
	"lemonde.fr", "npr.org", "macys.com", "51cto.com", "kijiji.ca", "newegg.com", "alnaharegypt.com", "flvto.biz",
	"namasha.com", "nhentai.net", "discover.com", "google.rs", "ilovepdf.com", "realtor.com", "sportmail.ru", "pexels.com",
	"feedly.com", "kayak.com", "leagueoflegends.com", "google.bg", "animeflv.net", "doorblog.jp", "ticketmaster.com",
	"people.com", "google.com.ly", "egy.best", "sputniknews.com", "nmisr.com", "elmundo.es", "lifewire.com",
	"media.tumblr.com", "ajel.sa", "geeksforgeeks.org", "hh.ru", "intel.com", "xiami.com", "naver.jp", "superuser.com",
	"southwest.com", "archiveofourown.org", "zukxd6fkxqn.com", "asana.com", "cqnews.net", "caliente.mx", "mercari.com",
	"hatenablog.com", "google.tm", "itv.com", "mathrubhumi.com", "bitbucket.org", "uploaded.net", "kizlarsoruyor.com",
	"cdiscount.com", "olx.com.br", "torrentz2.eu", "europa.eu", "scribol.com", "investopedia.com", "elzmannews.com",
	"hubspot.com", "3c.tmall.com", "mercadolibre.com.ve", "artstation.com", "rajasthan.gov.in", "jiameng.com",
	"thewhizmarketing.com", "t-online.de", "bola.com", "dailycaller.com", "labanquepostale.fr", "sportbible.com",
	"delta.com", "cricbuzz.com", "torrent9.blue", "wtoip.com", "gome.com.cn", "unsplash.com", "e-hentai.org",
	"indiamart.com", "verizonwireless.com", "myfreecams.com", "redfin.com", "livescore.com", "6.cn", "infusionsoft.com",
	"bicentenariobu.com", "animeyt.tv", "moneycontrol.com", "gazeta.ru", "wunderground.com", "zing.vn", "cdstm.cn",
	"youjizz.com", "google.tn", "justdial.com", "nordstrom.com", "tomsguide.com", "best2017games.com", "lowes.com",
	"dafont.com", "abc.net.au", "secureserver.net", "caixa.gov.br", "time.com", "baike.com", "bldaily.com", "lenovo.com",
	"vkuseraudio.net", "lequipe.fr", "beytoote.com", "nextlnk3.com", "pchome.com.tw", "pole-emploi.fr", "dribbble.com",
	"avgle.com", "bhphotovideo.com", "huawei.com", "farfetch.com", "xe.com", "firefoxchina.cn", "medianetto.com",
	"segmentfault.com", "agoda.com", "trustpilot.com", "dictionary.com", "meetup.com", "mileroticos.com", "yournewtab.com",
	"cbsnews.com", "asahi.com", "arxiv.org", "yaplakal.com", "uniqlo.tmall.com", "timeanddate.com", "espncricinfo.com",
	"ruliweb.com", "rutube.ru", "lefigaro.fr", "linkshrink.net", "www.gov.uk", "fmovies.se", "ninisite.com", "uptobox.com",
	"360doc.com", "viva.co.id", "mynet.com", "prom.ua", "auction.co.kr", "pcgamer.com", "exhentai.org", "dmm.com",
	"rapidgator.net", "prezi.com", "kp.ru", "japanpost.jp", "wayfair.com", "td.com", "buffstreamz.com", "avast.com",
	"mediaset.it", "gogoanime.se", "nianhuo.tmall.com", "nvzhuang.tmall.com", "issuu.com", "11st.co.kr", "cam4.com",
	"nature.com", "4pda.ru", "codepen.io", "donga.com", "ccm.net", "chip.de", "kissasian.ch", "1tv.ru", "apkpure.com",
	"getpocket.com", "dmv.org", "addthis.com", "eksisozluk.com", "ukr.net", "motherless.com", "tmz.com",
	"y8.com", "khanacademy.org", "tomshardware.com", "4shared.com", "canada.ca", "zol.com.cn", "alodokter.com",
	"taqviyat.com", "westernjournal.com", "chron.com", "thethao247.vn", "aa.com", "reundcwkqvctq.com", "vesti.ru",
	"trulia.com", "google.com.gt", "getcryptotab.com", "visualstudio.com", "mydrivers.com", "python.org",
	"bolasport.com", "state.gov", "nextlnk2.com", "mirror.co.uk", "usnews.com", "infobae.com", "creditkarma.com", "tf1.fr",
	"telegraf.com.ua", "zimuzu.tv", "hotnewhiphop.com", "academia.edu", "cimaclub.com", "gotomeeting.com", "costco.com",
	"gazzetta.it", "popcash.net", "4399.com", "knowyourmeme.com", "axisbank.co.in", "prnt.sc", "subito.it", "duba.com",
	"jqw.com", "mobile.de", "makemytrip.com", "badoo.com", "subject.tmall.com", "chess.com", "usaa.com", "united.com",
	"medianewpage.com", "admaimai.com", "ivi.ru", "sapo.pt", "skysports.com", "oschina.net", "sozcu.com.tr", "envato.com",
	"bestadbid.com", "harvard.edu", "hatena.ne.jp", "epfindia.gov.in", "gmanetwork.com", "thehill.com", "pinterest.co.uk",
	"bitmedianetwork.com", "duolingo.com", "elwatannews.com", "bankmellat.ir", "pastebin.com", "howtogeek.com", "xtube.com",
	"verizon.com", "worldstarhiphop.com", "miui.com", "fanfiction.net", "nikkei.com", "scol.com.cn", "americanas.com.br",
	"stanford.edu", "wikiwand.com", "mawdoo3.com", "sports.ru", "nba.com", "istockphoto.com", "google.lt", "wildberries.ru",
	"binance.com", "okcupid.com", "food.tmall.com", "google.iq", "onclicksuper.com", "turkiye.gov.tr",
	"lun.com", "pulzo.com", "sci-hub.tw", "google.hr", "bodybuilding.com", "olymptrade.com", "ebates.com", "gamib.com",
	"google.com.my", "dw.com", "crunchbase.com", "shopee.tw", "ruten.com.tw", "freejobalert.com", "concursolutions.com",
	"lentainform.com", "norton.com", "fishki.net", "qualtrics.com", "focus.de", "digitaltrends.com", "douyin.com", "ask.fm",
	"zcool.com.cn", "chegg.com", "n11.com", "v2ex.com", "francetvinfo.fr", "tempo.co", "chiphell.com", "chinadaily.com.cn",
	"redwap.me", "getbootstrap.com", "rutor.info", "motorsport.com", "marriott.com", "strava.com", "hootsuite.com",
	"urdupoint.com", "teamviewer.com", "ea.com", "reclameaqui.com.br", "zomato.com", "fazenda.gov.br", "ozon.ru",
	"ancestry.com", "ebay.fr", "myanmarload.com", "atlassian.com", "ecosia.org", "5ch.net", "nate.com",
	"tandfonline.com", "nbcnews.com",
}
