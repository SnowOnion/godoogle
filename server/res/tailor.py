#!/usr/bin/env python3 -i

"""
func[a any](bool, a, a) a
func[a any](bool, a, a) a
"""

import json
from pprint import pp

with open("../../../godoogle_archive/floyd.json", "r") as fi:
    j=json.load(fi)
ter_sig="func[T any](bool, T, T) T"
ternary=j[ter_sig]

j_only_ternary={ter_sig:j[ter_sig]}

with open("ter.json","w") as fo:
    json.dump(j_only_ternary,fo)
