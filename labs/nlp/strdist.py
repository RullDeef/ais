translite = {
    'а': 'a',
    'б': 'b',
    'в': 'v',
    'г': 'g',
    'д': 'd',
    'е': 'e',
    'ё': 'yo',
    'ж': 'j',
    'з': 'z',
    'и': 'i',
    'к': 'k',
    'л': 'l',
    'м': 'm',
    'н': 'n',
    'о': 'o',
    'п': 'p',
    'р': 'r',
    'с': 's',
    'т': 't',
    'у': 'u',
    'ф': 'f',
    'х': 'h',
    'ц': 'c',
    'ч': 'ch',
    'ш': 'sh',
    'щ': 'sh',
    'ъ': '',
    'ы': '',
    'ь': '',
    'э': 'e',
    'ю': 'u',
    'я': 'a',
    'w': 'v',
    
    'c': 'k',
    
    'e': 'a',
}


def transliterate(s: str) -> str:
    found = True
    while found:
        found = False
        for c, t in translite.items():
            if c in s:
                found = True
                s = s.replace(c, t)
                break
    return s


def str_distance(s1: str, s2: str) -> float:
    distances = [
        [0 for j in range(len(s2) + 1)]
        for i in range(len(s1) + 1)
    ]

    for t1 in range(len(s1) + 1):
        distances[t1][0] = t1

    for t2 in range(len(s2) + 1):
        distances[0][t2] = t2
        
    ins_cost = 1
    rem_cost = 1
    repl_cost = 0.6
    
    a, b, c = 0, 0, 0
    for t1 in range(1, len(s1) + 1):
        for t2 in range(1, len(s2) + 1):
            if s1[t1-1] == s2[t2-1]:
                distances[t1][t2] = distances[t1 - 1][t2 - 1]
            else:
                a = distances[t1][t2 - 1]
                b = distances[t1 - 1][t2]
                c = distances[t1 - 1][t2 - 1]
                if a <= b and a <= c:
                    distances[t1][t2] = a + ins_cost
                elif b <= a and b <= c:
                    distances[t1][t2] = b + rem_cost
                else:
                    distances[t1][t2] = c + repl_cost
    
    return distances[len(s1)][len(s2)]


def str_similarity(s1: str, s2: str) -> float:
    s1, s2 = s1.lower(), s2.lower()
    s1, s2 = transliterate(s1), transliterate(s2)
    dist = str_distance(s1, s2)
    return max(0, min(1 - dist / max(len(s1), len(s2)), 1))

    
if __name__ == "__main__":
    strs = (
        ('Экшен', 'Action'),
        ('Адвенчур', 'adventure'),
        ('Адвенчeр', 'adventure'),
        ('сёнен', 'seinen'),
        ('Omega', 'пепега'),
    )
    
    for s1, s2 in strs:
        perc = 100 * str_similarity(s1, s2)
        print(f'sim({s1}, {s2}) = {perc:.0f}%')
