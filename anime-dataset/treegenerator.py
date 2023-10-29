from typing import Callable
import math
import random


class PropTree:
    def __init__(self, values: list[str], *, links: list[int]=None):
        self._N = len(values)
        self._values = values[:]
        self._links = [None] * self._N # every value points to its parent or None, if has no parent
        # distances matrix (must be rebuild when changing self._links)
        self._distances = [[(0 if i == j else self._N) for j in range(self._N)] for i in range(self._N)]
        self._height = 1
        if links is not None:
            self._links = links
            self._recalc_distances()


    # returns distance in steps between values in current configuration
    def distance(self, value1: str, value2: str) -> int:
        i = self.get_value_index(value1)
        j = self.get_value_index(value2)
        return self._distances[i][j]
    
    
    def leaf_height(self, value: str) -> int:
        h = 0
        id = self._values.index(value)
        while id is not None:
            id = self._links[id]
            h += 1
        return h


    def height(self) -> int:
        return self._height
    
    
    def show(self):
        def show_subtree(par_id: int, tabs=''):
            print(tabs + self._values[par_id])
            for i in range(self._N):
                if self._links[i] == par_id:
                    show_subtree(i, tabs + '|-')
        print('PropTree links:', self._links)
        for i in range(self._N):
            if self._links[i] is None:
                show_subtree(i)


    # updates value's parent to new_parent
    def reparent(self, value: str, new_parent: str):
        id = self.get_value_index(value)
        pid = self.get_value_index(new_parent)
        old_parent = self._links[id]
        self._links[id] = pid
        if self._detect_cycle():
            self._links[id] = old_parent
            raise ValueError('Cannot reparent: cycle in prop tree detected')
        self._recalc_distances()


    def get_value_index(self, value: str) -> int:
        try:
            return self._values.index(value)
        except ValueError:
            raise ValueError(f'value "{value}" not in prop tree')


    def _detect_cycle(self) -> bool:
        parents = list(range(self._N))
        for i in range(self._N):
            was_change = False
            for v in range(self._N):
                if self._links[parents[v]] is not None:
                    parents[v] = self._links[parents[v]]
                    was_change = True
            if not was_change:
                return False
        return True


    def _recalc_distances(self):
        def build_path(id: int) -> set[int]:
            par = self._links[id]
            return set((id,)) if par is None else build_path(par).union((id,))
        val_paths = [build_path(i) for i in range(self._N)]
        self._height = 1
        for i in range(self._N):
            self._height = max(self._height, len(val_paths[i]))
            for j in range(i + 1, self._N):
                d = len(val_paths[i].symmetric_difference(val_paths[j]))
                self._distances[i][j] = d
                self._distances[j][i] = d


class PropTreeProxy:
    def __init__(prop):
        pass


def load_from_str(string: str) -> PropTree:
    strings = string.split('\n')
    values = [s.lstrip('|-') for s in strings]
    links = [None] * len(values)
    def recursive_loader(par_id: int, k: int):
        if '|-' * (k // 2) != strings[0][:k]:
            return False
        id = values.index(strings[0][k:])
        links[id] = par_id
        strings.pop(0)
        while len(strings) > 0 and strings[0].startswith('|-' * (k // 2)):
            if not recursive_loader(id, k + 2):
                break
        return True
    while len(strings) != 0:
        recursive_loader(None, 0)
    return PropTree(values, links=links)


class MutatablePropTree(PropTree):
    def __init__(self, values: list[str], *, links: list[int]=None):
        if links is not None:
            super().__init__(values, links=links)
        else:
            super().__init__(values)


    # performs random mutation    
    def mutate(self, *, rate=0.1):
        if random.random() >= rate:
            return
        while True:
            try:
                val, parent = random.choices(self._values, k=2)
                self.reparent(val, parent)
                break
            except ValueError:
                pass
    
    
    def clone(self):
        copy = MutatablePropTree(self._values)
        for i in range(self._N):
            par = self._links[i]
            if par is not None:
                copy.reparent(self._values[i], self._values[par])
        return copy


class PropTreeGeneticGenerator:
    def __init__(self, values: list[str], links: list[int], metric: Callable[[PropTree], float], *, population: int=100):
        self._population = [MutatablePropTree(values, links=links) for i in range(population)]
        self._metric = metric
        self._mutation_rate = 0.8
        self._survival_rate = 0.7


    def best_ent(self) -> PropTree:
        return self._population[0]


    def metric_mean(self):
        return sum(self._metric(p) for p in self._population) / len(self._population)


    def epoch(self):
        for p in self._population:
            p.mutate(rate=self._mutation_rate)
        self._population.sort(key=self._metric)
        pop_size = len(self._population)
        n = math.ceil(self._survival_rate * pop_size)
        self._population[n:] = []
        while len(self._population) < pop_size:
            i = random.randint(0, n-1)
            p = self._population[i]
            self._population.insert(i, p.clone())


def load_best_genres_tree() -> PropTree:
    return load_from_str('''Slice of Life
|-Comedy
|-Drama
|-Music
|-Parody
|-School
Shoujo
|-Josei
|-Romance
Shounen
|-Adventure
|-Game
|-Sports
|-|-Cars
|-Action
|-|-Martial Arts
|-|-Historical
|-|-|-Samurai
|-|-Police
|-|-|-Military
Seinen
|-Psychological
Sci-Fi
|-Mecha
|-Space
Mystery
|-Thriller
|-|-Horror
|-Super Power
|-Supernatural
|-|-Demons
|-|-Vampire
|-Fantasy
|-|-Magic''')


def compute_similarity(propTree: PropTree, genres1: list[str], genres2: list[str]) -> float:
    gs1, gs2 = set(genres1), set(genres2)
    jac = len(gs1.intersection(gs2)) / (len(gs1.union(gs2)))
    
    gs1, gs2 = gs1.difference(gs2), gs2.difference(gs1)
    gs1, gs2 = list(gs1), list(gs2)

    sim = 0
    if len(gs1) != 0 and len(gs2) != 0:
        dist_mat = [[propTree.distance(g_i, g_j) for g_j in gs1] for g_i in gs2]
        max_dist = 1
        for row in dist_mat:
            for el in row:
                max_dist = max(max_dist, el)
        for i in range(len(dist_mat)):
            for j in range(len(dist_mat[i])):
                sim += 1 - dist_mat[i][j] / max_dist
        sim /= len(genres1) * len(genres2)
    if len(gs1) == 0 and len(gs2) != 0:
        sim2to1 = [[propTree.distance(g_i, g_j) for g_j in gs2] for g_i in genres1]
        max_dist = 1
        for row in sim2to1:
            for el in row:
                max_dist = max(max_dist, el)
        for i in range(len(sim2to1)):
            for j in range(len(sim2to1[i])):
                sim += 1 - sim2to1[i][j] / max_dist
        sim /= len(gs2)
    if len(gs1) != 0 and len(gs2) == 0:
        sim1to2 = [[propTree.distance(g_i, g_j) for g_j in gs1] for g_i in genres2]
        max_dist = 1
        for row in sim1to2:
            for el in row:
                max_dist = max(max_dist, el)
        for i in range(len(sim1to2)):
            for j in range(len(sim1to2[i])):
                sim += 1 - sim1to2[i][j] / max_dist
        sim /= len(gs1)
    return min(1, jac + sim / 6)


if __name__ == '__main__':
    import pandas as pd
    import numpy as np
    
    dataset = pd.read_csv('animedb/data/anime.csv')
    dataset = dataset.drop(columns=[f'Score-{i}' for i in range(1, 11)])
    genres = [
        'Psychological',
        'Action',
        'Shounen',
        'Supernatural',
        # 'Dementia',
        'Sports',
        'Martial Arts',
        'Historical',
        'Demons',
        'Josei',
        'Space',
        # 'Harem',
        # 'Kids',
        # 'Ecchi',
        'Mystery',
        'Vampire',
        # 'Shounen Ai',
        # 'Yuri',
        'Cars',
        'Super Power',
        'Seinen',
        'Sci-Fi',
        'Magic',
        'Parody',
        'Thriller',
        'Music',
        'Game',
        'Fantasy',
        'Adventure',
        'Romance',
        'Police',
        'Drama',
        # 'Hentai',
        'Samurai',
        'School',
        'Comedy',
        'Shoujo',
        # 'Yaoi',
        'Military',
        'Horror',
        'Slice of Life',
        # 'Shoujo Ai',
        'Mecha'
    ]
    
    filtered_genres = ['Dementia', 'Harem', 'Kids', 'Ecchi', 'Shounen Ai', 'Yuri', 'Hentai', 'Yaoi', 'Shoujo Ai']
    def genre_filter(genre_lst: str) -> bool:
        for g in genre_lst.split(', '):
            if g in filtered_genres:
                return False
        return True

    dataset = dataset.loc[dataset['Genres'].apply(genre_filter)]
    
    titles = (
        'Fullmetal Alchemist:Brotherhood',
        'Attack on Titan',
        'Jujutsu Kaisen',
        
        'Bleach',
        'Naruto',
        'One Piece',
        
        'Haikyu!!',
        'Kuroko\'s Basketball',
        'Run with the Wind',

        'Violet Evergarden',
        'Your Lie in April',
        'Fruits Basket',
        
        'Made in Abyss',
        'Steins;Gate',
        'The Promised Neverland',
        
        'Gintama',
        'Great Teacher Onizuka',
        'Nichijou - My Ordinary Life',
    )

    featured = dataset[dataset['English name'].isin(titles)]
    featured = featured.loc[featured['English name'].apply(lambda s: titles.index(s)).sort_values().index]
    
    def metric_fn(propTree: PropTree) -> float:
        res = 0
        for i in range(len(titles)):
            genres_i = featured.iloc[i]['Genres'].split(', ')
            for j in range(i + 1, len(titles)):
                genres_j = featured.iloc[j]['Genres'].split(', ')
                sim = compute_similarity(propTree, genres_i, genres_j)
                want = 1 if i // 3 == j // 3 else 0
                res += abs(sim - want)
        return res + 3 * propTree.height()

    best_tree = load_best_genres_tree()
    
    print('metric:', metric_fn(best_tree))
    
    print('test sim:', compute_similarity(best_tree,
        ['Seinen', 'Sci-Fi', 'Fantasy'],
        ['Fantasy', 'Super Power', 'Magic'],
    ))
    exit()

    default_links = best_tree._links[:]
    generator = PropTreeGeneticGenerator(best_tree._values, default_links, metric_fn)
    
    print('initial metric mean:', generator.metric_mean())
    for i in range(1000):
        generator.epoch()
        print(f'epoch #{i} metric mean:', generator.metric_mean())
        generator.best_ent().show()
    
    
    # best_tree.show()

'''
Best so far:
epoch #174 metric mean: 55.174615440115424
PropTree links: [22, 23, 30, 21, 22, 4, 12, 0, 5, 33, 14, 17, 10, 10, 22, 3, 22, 22, 16, 33, 10, 4, None, 24, 17, 16, 25, 12, 29, 7, 4, 9, 12, 22]
Adventure
|-Psychological
|-|-Demons
|-|-|-Shoujo
|-|-|-|-Comedy
|-Sports
|-|-Martial Arts
|-|-|-Josei
|-|-Fantasy
|-|-|-Supernatural
|-|-|-|-Sci-Fi
|-|-Military
|-|-|-Shounen
|-Seinen
|-|-Mystery
|-|-|-Cars
|-|-|-|-Historical
|-|-|-|-School
|-|-|-|-Slice of Life
|-|-|-Super Power
|-|-|-Game
|-Magic
|-|-Thriller
|-|-Drama
|-|-|-Samurai
|-Parody
|-|-Vampire
|-|-Police
|-|-|-Romance
|-|-|-|-Action
|-Mecha
|-|-Space
|-|-|-Horror
|-|-Music
'''
