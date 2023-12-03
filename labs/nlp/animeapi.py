import requests
from typing import Optional
from strdist import str_similarity

""" Example DTO:

{
    "Id": 11013,
    "Title": "Inu x Boku SS",
    "ImageURL": "https://cdn.myanimelist.net/images/anime/12/35893.jpg",
    "Type": "TV",
    "Source": "Manga",
    "Studio": "David Production",
    "Genres": [
        "Comedy",
        "Supernatural",
        "Romance",
        "Shounen"
    ],
    "Duration": 24,
    "Episodes": 12,
    "Year": 2012,
    "AiredFrom": "2012-01-13T00:00:00Z",
    "AiredTo":"2012-03-30T00:00:00Z"
}
"""

class AnimeDTO:
    def __init__(self, *, Id: int, Title: str, ImageURL: str, Type: str, Source: str, Studio: str, Genres: list[str], Duration: int, Episodes: int, Year: int, AiredFrom: str, AiredTo: str):
        self.id = Id
        self.title = Title
        self.image_url = ImageURL
        self.type = Type
        self.source = Source
        self.studio = Studio
        self.genres = Genres
        self.duration = Duration
        self.episodes = Episodes
        self.year = Year
        self.aired_from = AiredFrom
        self.aired_to = AiredTo
    
    def __str__(self) -> str:
        return f'AnimeDTO(id={self.id}, title="{self.title}")'
    
    def __repr__(self) -> str:
        return str(self)


class ApiServer:
    def __init__(self, host: str, port: int):
        self.__base_url = f'http://{host}:{port}/api'
    
    def get_animes(self) -> list[AnimeDTO]:
        resp = requests.get(f'{self.__base_url}/animes')
        if not resp.ok:
            print(resp.reason)
        return [AnimeDTO(**values) for values in resp.json()]
    
    def get_filtered_animes(self, page: int) -> list[AnimeDTO]:
        resp = requests.get(f'{self.__base_url}/animes', {'page': page})
        if not resp.ok:
            print(resp.reason)
        return [AnimeDTO(**values) for values in resp.json()]

    def like_anime(self, anime_id: int):
        resp = requests.get(f'{self.__base_url}/animes/{anime_id}', {'mark': 'fav'})
        if not resp.ok:
            print(resp.reason)
    
    def dislike_anime(self, anime_id: int):
        resp = requests.get(f'{self.__base_url}/animes/{anime_id}', {'mark': 'unfav'})
        if not resp.ok:
            print(resp.reason)
    
    def clear_preference(self, anime_id: int):
        resp = requests.get(f'{self.__base_url}/animes/{anime_id}', {'mark': 'clear'})
        if not resp.ok:
            print(resp.reason)

    def get_recomendations(self, page: int) -> list[AnimeDTO]:
        resp = requests.get(f'{self.__base_url}/recomendations', {'page': page})
        if not resp.ok:
            print(resp.reason)
        data = resp.json()
        if data is None:
            return []
        return [AnimeDTO(**values) for values in data]


class AnimeService:
    def __init__(self, apiServer: ApiServer):
        self.__api = apiServer
        self.__animes = apiServer.get_animes()
    
    def find_anime_exact(self, query: str) -> Optional[AnimeDTO]:
        min_assurance = 0.7
        animes = [(str_similarity(a.title, query), a) for a in self.__animes]
        assurance, anime = max(animes, key=lambda a: a[0])
        animes = list(filter(lambda s: s[0] > 0.7, animes))
        if assurance < min_assurance:
            return None
        return anime
    
    def find_anime_fuzzy(self, query: str) -> list[AnimeDTO]:
        min_assurance = 0.5
        animes = [(str_similarity(a.title, query), a) for a in self.__animes]
        animes = list(filter(lambda s: s[0] > min_assurance, animes))
        animes.sort(key=lambda a: a[0], reverse=True)
        return [a[1] for a in animes]
    
    def like_anime(self, anime_id: int):
        self.__api.like_anime(anime_id)
    
    def dislike_anime(self, anime_id: int):
        self.__api.dislike_anime(anime_id)
    
    def get_recomendations(self) -> list[AnimeDTO]:
        recs = self.__api.get_recomendations(1)
        if len(recs) == 0:
            recs = self.__api.get_filtered_animes(1)
        return recs
    
    def get_recomendations_str(self) -> str:
        return "\n".join([f'{i+1}) [ID#{a.id}] {a.title}'
                          for i, a in enumerate(self.get_recomendations())])


if __name__ == "__main__":
    api = ApiServer('localhost', 8080)
    service = AnimeService(api)
    
    animes = service.get_recomendations()
    print(animes)
    
    query = 'Блич'
    anime = service.find_anime_exact(query)
    print(anime)
