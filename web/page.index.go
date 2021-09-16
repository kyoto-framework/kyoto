package main

import (
	"html/template"

	"github.com/yuriizinets/go-ssc"
)

type PageIndex struct {
	CMeta      ssc.Component
	Navbar     ssc.Component
	Footer     ssc.Component
	Features   ssc.Component
	Sponsors   ssc.Component
	Statistics ssc.Component
	Promo      ssc.Component
	DevPromo   ssc.Component
	Details    ssc.Component
	About      ssc.Component
	FAQ        ssc.Component
}

func (p *PageIndex) Template() *template.Template {
	return template.Must(template.New("page.index.html").Funcs(tfuncs()).ParseGlob("*.html"))
}

func (p *PageIndex) Meta() ssc.Meta {
	return ssc.Meta{
		Title: "Home - SSC",
	}
}

func (p *PageIndex) Init() {
	p.CMeta = ssc.RegC(p, &ComponentMeta{})
	p.Navbar = ssc.RegC(p, &ComponentNavbar{
		Left: []ComponentNavbarHref{
			{
				Title: "Documentation",
				Href:  "/docs/",
			},
			{
				Title: "Contact",
				Href:  "mailto:yurii.zinets@icloud.com",
			},
		},
		Right: []ComponentNavbarHref{
			{
				Image: template.HTML(`
				<svg class="h-6 w-6 dark:text-white" viewBox="0 0 256 250" fill="currentColor" stroke="currentColor"  version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="xMidYMid">
					<g>
						<path fill="currentColor" d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z" fill="#161614"></path>
					</g>
				</svg>
				`),
				Href: "https://github.com/yuriizinets/go-ssc",
			},
		},
	})
	p.Footer = ssc.RegC(p, &ComponentFooter{
		Left: []ComponentFooterHref{
			{
				Title: "Documentation",
				Href:  "/docs/",
			},
			{
				Title: "Contact",
				Href:  "mailto:yurii.zinets@icloud.com",
			},
		},
		Center: []ComponentFooterHref{
			{
				Image: template.HTML(`
				<svg class="h-6 w-6 dark:text-white" viewBox="0 0 256 250" fill="currentColor" stroke="currentColor"  version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" preserveAspectRatio="xMidYMid">
					<g>
						<path fill="currentColor" d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z" fill="#161614"></path>
					</g>
				</svg>
				`),
				Href: "https://github.com/yuriizinets/go-ssc",
			},
		},
		Copyright: "© 2021-present Created by Yurii Zinets",
	})
	p.Features = ssc.RegC(p, &ComponentBlockFeatures{
		Entries: []ComponentBlockFeaturesEntry{
			{
				Image:       template.HTML(`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z"></path></svg>`),
				Title:       "Functional",
				Description: "Component approach, async operations, component server side methods and more",
			},
			{
				Image:       template.HTML(`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path></svg>`),
				Title:       "Customizable",
				Description: "Feel free to use context feature, define own handlers with context setters",
			},
			{
				Image:       template.HTML(`<svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>`),
				Title:       "Pure Design",
				Description: "You're free to use any framework as base for your web project, net/http is not a requirement",
			},
		},
	})
	p.Sponsors = ssc.RegC(p, &ComponentSponsors{
		Title: "Our sponsors",
		Entries: []ComponentSponsorsEntry{
			{
				Photo: "/static/img/sponsors/brokerone.svg",
				Href:  "https://mybrokerone.com",
			},
		},
	})
	p.Promo = ssc.RegC(p, &ComponentContent{
		Title:       "SSR-first, fast, reliable",
		Description: "An HTML render engine concept that brings frontend-like components experience to the server side with native html/template on steroids. Ideal fit for SEO strict projects. Fast, on-demand HTML, minimal JS payload.",
		Photo:       "/static/img/demo.png",
		LinkTitle:   "Check Cons & Pros",
		LinkHref:    "/docs/",
	})
	p.DevPromo = ssc.RegC(p, &ComponentContent{
		Title:       "Developer friendly",
		Description: "This library was created to solve developer problems first. Use full power of Go, build your parallel development process around components system, deliver your webpages fast! Feel free to use landing page of this library as starter project.",
		Photo:       "/static/img/code.jpg",
		Reverse:     true,
		LinkTitle:   "Check Documentation",
		LinkHref:    "/docs/",
	})

	p.Details = ssc.RegC(p, &ComponentBlockFeaturesMore{
		Title:    "core features",
		Subtitle: "A better way to build your frontend",
		Entries: []ComponentBlockFeaturesEntry{
			{
				Image:       template.HTML(`<svg class="w-10 h-10 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path></svg>`),
				Title:       "html/template",
				Description: "No custom template engines, just use built-in one!",
			},
			{
				Image:       template.HTML(`<svg class="w-10 h-10 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>`),
				Title:       "Asynchronous operations",
				Description: "Hassle-free asynchronous methods. It's enough to define a method, that's all!",
			},
			{
				Image:       template.HTML(`<svg class="w-10 h-10 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>`),
				Title:       "Context",
				Description: "Use page-level context without worring about concurrency and mutex",
			},
			{
				Image:       template.HTML(`<svg class="w-10 h-10 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3 3m0 0l-3-3m3 3V4"></path></svg>`),
				Title:       "Server Side Actions",
				Description: "Use server defined component methods, instead of including logic in JS payload",
			},
			{
				Image:       template.HTML(`<svg class="w-10 h-10 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"></path></svg>`),
				Title:       "Cross-component communication",
				Description: "Trigger actions, defined in different components",
			},
		},
	})
	p.Statistics = ssc.RegC(p, &ComponentBlockStatistics{
		Entries: []ComponentBlockStatisticsEntry{
			{
				Title: "stars",
				Count: 63,
			},
			{
				Title: "contributors",
				Count: 3,
			},
			{
				Title: "projects",
				Count: 2,
			},
			{
				Title: "sponsors",
				Count: 1,
			},
		},
	})
	p.About = ssc.RegC(p, &ComponentBlockFAQ{
		Title: "About",
		Entries: []ComponentBlockFAQEntry{
			{
				Question: "Motivation",
				Answer:   "Main motivation is to reduce usage of popular SPA/PWA frameworks where it's not needed because it adds a lot of complexity and overhead. There is no reason to bring significant runtime, VirtualDOM, and Webpack into the project with minimal dynamic frontend behavior. This project proves the possibility of keeping most of the logic on the server's side.",
			},
			{
				Question: "What problems does it solve?",
				Answer: template.HTML(`
					While developing the website's frontend with traditional Go handlers and templates, I discovered some of the downsides of this approach:
					<ul class="py-2 list-disc list-inside">
						<li>With plain html/template you're starting to repeat yourself. It's harder to define reusable parts.</li>
						<li>You must repeat DTO calls for each page, where you're using reusable parts.</li>
						<li>With Go's routines approach it's hard to make async-like DTO calls in the handlers.</li>
						<li>For dynamic things, you still need to use JS and client-side DOM modification.</li>
					</ul>
					Complexity is much higher when all of them get combined. This engine tries to bring components and async experience to the traditional server-side rendering.
				`),
			},
			{
				Question: "Zen",
				Answer: template.HTML(`
				For contributors:
				<ul class="py-2 list-disc list-inside">
					<li>Don't replace Go features that exist already</li>
					<li>Don't do work that's already done</li>
					<li>Don't force developers to use a specific solution (Gin/Chi/GORM/sqlx/etc). Let them choose</li>
					<li>Rely on the server to do the rendering, minimum JS specifics or client-side only behavior</li>
					<li>KISS</li>
				</ul>
				`),
			},
		},
	})
	p.FAQ = ssc.RegC(p, &ComponentBlockFAQ{
		Title: "Frequently-asked questions",
		Entries: []ComponentBlockFAQEntry{
			{
				Question: "We already have ongoing projects. Should we rewrite everything?",
				Answer:   "No! Just simply define handlers, and use this library for creating new pages.",
			},
			{
				Question: "Is this library ready for production?",
				Answer:   "No, we don't recommend to use it in production (even if we are using so).",
			},
			{
				Question: "Seems like SSA feature is slow, is it a bug?",
				Answer:   "No, you need to check documentation first to explore cons & pros of SSA. Use SSA carefully",
			},
		},
	})
}
