using System;
using System.Collections.Generic;
using System.IO;
using System.Threading;

namespace ScopusCrawler.Processing
{
    public static class KeysStorage
    {
        private static List<string> keys;

        private static void Initialize()
        {
            string[] readKeys = File.ReadAllLines("data\\keys.txt");
            var result = new List<string>();
            foreach (var key in readKeys)
            {
                result.Add(key);
            }
            keys = result;
        }

        public static string GetKey()
        {
            if (keys == null)
                Initialize();
            var random = new Random();
            var sleepTime = random.Next(100);
            Thread.Sleep(sleepTime);
            var index = random.Next(keys.Count - 1);
            return keys[index];
        }
    }
}
