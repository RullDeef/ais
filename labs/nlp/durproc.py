import math
from strdist import str_similarity

durations = (
    ("очень короткий", 0, 10),
    ("короткий", 10, 35),
    ("недолгий", 35, 55),
    ("не очень долгий", 55, 81),
    ("не очень короткий", 81, 96),
    ("некороткий", 96, 105),
    ("долгий", 105, 130),
    ("очень долгий", 130, math.inf),
)

def is_duration(s: str) -> [float, tuple[str, float, float]]:
    min_similarity = 0.7
    sim_map = [(str_similarity(s, dd[0]), dd) for dd in durations]
    res = max(sim_map, key=lambda s: s[0] + 0.1 * len(s[1][0].split()))
    if res[0] <= min_similarity:
        return 0, ''
    return res
