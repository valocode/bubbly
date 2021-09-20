package adapter

code_scan[scan] {
	scan := {"tool": "meta-doubleopen"}
}

component := [comp |
	pkg := input[_].packages[_]

	# check that pkg has the key externalRefs and versionInfo, otherwise skip it
	pkg.externalRefs
	pkg.versionInfo

	comp := {
		"name": pkg.name,
		"version": pkg.versionInfo,
		"description": pkg.description,
		"url": pkg.homepage,
		"vulnerabilities": [vuln |
			# get the list of fixed (patched) CVEs
			fixedCVE := split(trim_prefix(pkg.sourceInfo, "CVEs fixed: "), " ")[_]

			vuln := {
				"vid": fixedCVE,
				"patch": {"note": "patch detected by meta-doubleopen"},
			}
		],
		"licenses": [lic |
			not pkg.licenseDeclared == "NOASSERTION"

			# split the licenses by their ANDs and ORs, which come in two flavours
			licList := split(split(split(split(pkg.licenseDeclared, "|")[_], "OR")[_], "&")[_], "AND")

			# trim each license to remove whitespace and any brackets
			lic_id := trim(trim(trim_space(licList[_]), "("), ")")

			# TODO: what if licenses did not map? We should not just silently ignore it
			# create the license from the processed license IDs
			lic := {"license_id": spdx_license[lic_id]}
		],
	}
]

# Source for SPDX license mapping
# http://git.yoctoproject.org/cgit.cgi/poky/tree/meta/conf/licenses.conf
spdx_license := {
	# AGPL variations
	"AGPL-3": "AGPL-3.0-only",
	"AGPL-3+": "AGPL-3.0-or-later",
	"AGPLv3": "AGPL-3.0-only",
	"AGPLv3+": "AGPL-3.0-or-later",
	"AGPLv3.0": "AGPL-3.0-only",
	"AGPLv3.0+": "AGPL-3.0-or-later",
	"AGPL-3.0": "AGPL-3.0-only",
	"AGPL-3.0+": "AGPL-3.0-or-later",
	# BSD variations
	"BSD-0-Clause": "0BSD",
	# GPL variations
	"GPL-1": "GPL-1.0-only",
	"GPL-1+": "GPL-1.0-or-later",
	"GPLv1": "GPL-1.0-only",
	"GPLv1+": "GPL-1.0-or-later",
	"GPLv1.0": "GPL-1.0-only",
	"GPLv1.0+": "GPL-1.0-or-later",
	"GPL-1.0": "GPL-1.0-only",
	"GPL-1.0+": "GPL-1.0-or-later",
	"GPL-2": "GPL-2.0-only",
	"GPL-2+": "GPL-2.0-or-later",
	"GPLv2": "GPL-2.0-only",
	"GPLv2+": "GPL-2.0-or-later",
	"GPLv2.0": "GPL-2.0-only",
	"GPLv2.0+": "GPL-2.0-or-later",
	"GPL-2.0": "GPL-2.0-only",
	"GPL-2.0+": "GPL-2.0-or-later",
	"GPL-3": "GPL-3.0-only",
	"GPL-3+": "GPL-3.0-or-later",
	"GPLv3": "GPL-3.0-only",
	"GPLv3+": "GPL-3.0-or-later",
	"GPLv3.0": "GPL-3.0-only",
	"GPLv3.0+": "GPL-3.0-or-later",
	"GPL-3.0": "GPL-3.0-only",
	"GPL-3.0+": "GPL-3.0-or-later",
	# LGPL variations
	"LGPLv2": "LGPL-2.0-only",
	"LGPLv2+": "LGPL-2.0-or-later",
	"LGPLv2.0": "LGPL-2.0-only",
	"LGPLv2.0+": "LGPL-2.0-or-later",
	"LGPL-2.0": "LGPL-2.0-only",
	"LGPL-2.0+": "LGPL-2.0-or-later",
	"LGPL2.1": "LGPL-2.1-only",
	"LGPL2.1+": "LGPL-2.1-or-later",
	"LGPLv2.1": "LGPL-2.1-only",
	"LGPLv2.1+": "LGPL-2.1-or-later",
	"LGPL-2.1": "LGPL-2.1-only",
	"LGPL-2.1+": "LGPL-2.1-or-later",
	"LGPLv3": "LGPL-3.0-only",
	"LGPLv3+": "LGPL-3.0-or-later",
	"LGPL-3.0": "LGPL-3.0-only",
	"LGPL-3.0+": "LGPL-3.0-or-later",
	# MPL variations
	"MPL-1": "MPL-1.0",
	"MPLv1": "MPL-1.0",
	"MPLv1.1": "MPL-1.1",
	"MPLv2": "MPL-2.0",
	# MIT variations
	"MIT-X": "MIT",
	"MIT-style": "MIT",
	# Openssl variations
	"openssl": "OpenSSL",
	# PSF variations
	"PSF": "PSF-2.0",
	"PSFv2": "PSF-2.0",
	# Python variations
	"Python-2": "Python-2.0",
	# Apache variations
	"Apachev2": "Apache-2.0",
	"Apache-2": "Apache-2.0",
	# Artistic variations
	"Artisticv1": "Artistic-1.0",
	"Artistic-1": "Artistic-1.0",
	# Academic variations
	"AFL-2": "AFL-2.0",
	"AFL-1": "AFL-1.2",
	"AFLv2": "AFL-2.0",
	"AFLv1": "AFL-1.2",
	# CDDL variations
	"CDDLv1": "CDDL-1.0",
	"CDDL-1": "CDDL-1.0",
	# Other variations
	"EPLv1.0": "EPL-1.0",
	"FreeType": "FTL",
	"Nauman": "Naumen",
	"tcl": "TCL",
	"vim": "Vim",
	# Silicon Graphics variations
	"SGIv1": "SGI-1",
}
