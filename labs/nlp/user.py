class PreferenceSet:
    def __init__(self):
        self.__likes = []
        self.__dislikes = []
    
    @property
    def likes(self):
        return self.__likes[:]
    
    @property
    def dislikes(self):
        return self.__dislikes[:]
    
    def like(self, item):
        if item not in self.__likes:
            self.__likes.append(item)
            self.__dislikes.remove(item)
    
    def dislike(self, item):
        if item not in self.__dislikes:
            self.__dislikes.append(item)
            self.__likes.remove(item)
    
    def like_many(self, items):
        for item in items:
            self.like(item)
    
    def dislike_many(self, items):
        for item in items:
            self.dislike(item)


class UserContext:
    def __init__(self):
        self.__anime_prefs = PreferenceSet()
        self.__genre_prefs = PreferenceSet()
        self.__filters = []

    @property
    def anime_prefs(self):
        return self.__anime_prefs

    @property
    def genre_prefs(self):
        return self.__genre_prefs
    
    def like_genres(self, genres: list[str]):
        self.__genre_prefs.like_many(genres)
    
    def dislike_genres(self, genres: list[str]):
        self.__genre_prefs.dislike_many(genres)
