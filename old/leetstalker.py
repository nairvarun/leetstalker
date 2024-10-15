import fire
import requests
import json
import asyncio
import aiohttp
import argparse


async def get_user_data(username: str):
  url: str = "https://leetcode.com/graphql/"
  json_data: dict = {
    "query": """
      query userPublicProfileAndProblemsSolved($username: String!) {
        userContestRanking(username: $username) {
          rating
          globalRanking
          topPercentage
        }
        matchedUser(username: $username) {
          profile {
            ranking
          }
          submitStatsGlobal {
            acSubmissionNum {
              difficulty
              count
            }
          }
        }
      }
    """,
    "variables": {
      "username": f"{username}",
    },
  }

  async with aiohttp.ClientSession() as session:
    async with session.post(url, json=json_data) as response:
      # handle errors
      # {'errors': [{'message': 'User matching query does not exist.', 'locations': [{'line': 3, 'column': 17}], 'path': ['userContestRanking']}, {'message': 'That user does not exist.', 'locations': [{'line': 8, 'column': 17}], 'path': ['matchedUser'], 'extensions': {'handled': True}}], 'data': {'userContestRanking': None, 'matchedUser': None}}
      if response.status == 200:
        data = await response.json()
        print(data)
      else:
        print(f"Failed with status code: {response.status}")


class LeetStalker:
  async def __call__(self):
    await get_user_data("nairvarun")

  def get(self, username: str):
    print("get")

  def add(self, username: str):
    print(username)

  def remove(self, username: str):
    print(username)

  def generate_config(self):
    print("gen conf")


if __name__ == "__main__":
  leetstalker = LeetStalker()
  fire.Fire(leetstalker)
