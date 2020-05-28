using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;

namespace ContestOzonGolangSchool.F
{
    public class Program
    {
        const string _inputFileName = "input.txt";

        public static void Main(string[] args)
        {
#if DEBUG
            var sw = new System.Diagnostics.Stopwatch();
            sw.Start();
            Console.WriteLine("Total memory 1: " + GC.GetTotalMemory(false));
#endif
            var result = 0;
            using (var inputStream = File.Exists(_inputFileName) ? File.OpenRead(_inputFileName) : Console.OpenStandardInput())
            using (var inputReader = new StreamReader(inputStream))
            {
                var sum = int.Parse(inputReader.ReadLine());
#if DEBUG
                Console.WriteLine("Total memory 2: " + GC.GetTotalMemory(false));
#endif
                result = SumExistsReadByWord(sum, inputReader);
            }

            File.WriteAllText("output.txt", result.ToString());

#if DEBUG
            Console.WriteLine("Total memory 3: " + GC.GetTotalMemory(false));
            Console.WriteLine("Elapsed: " + sw.ElapsedMilliseconds);
            Console.WriteLine("Result: " + result);
            Console.ReadLine();
#endif
        }

        private static string GetWord(StreamReader reader)
        {
            string res = null;
            while (!reader.EndOfStream)
            {
                var symbol = (char)reader.Read();
                if (char.IsWhiteSpace(symbol))
                {
                    return res;
                }
                res += symbol;
            }
            return res;
        }

        private static int SumExistsReadByWord(int sum, StreamReader reader)
        {
            var hashMap = new HashSet<int>();
            string word;
            while ((word = GetWord(reader)) != null)
            {
                var digit = int.Parse(word);
                if (digit >= sum)
                {
                    continue;
                }
                if (hashMap.Contains(digit))
                {
                    return 1;
                }
                hashMap.Add(sum - digit);
            }
            return 0;
        }

        private static int SumExistsSubstring(int sum, StreamReader reader)
        {
            var hashMap = new HashSet<int>();
            int idx;
            var digitsStr = reader.ReadLine();
            for (var prev = 0; (idx = digitsStr.IndexOf(" ", prev, StringComparison.Ordinal)) != -1; prev = idx + 1)
            {
                var digit = int.Parse(digitsStr.Substring(prev, idx - prev));
                if (digit >= sum)
                {
                    continue;
                }
                if (hashMap.Contains(digit))
                {
                    return 1;
                }
                hashMap.Add(sum - digit);
            }
            return 0;
        }

        private static int SumExistsArray(int sum, StreamReader reader)
        {
            var hashMap = new HashSet<int>();
            var digits = reader.ReadLine().Split(' ').Select(e => int.Parse(e)).ToArray();
            foreach (var digit in digits)
            {
                if (digit >= sum)
                {
                    continue;
                }
                if (hashMap.Contains(digit))
                {
                    return 1;
                }
                hashMap.Add(sum - digit);
            }
            return 0;
        }
    }


    class SimpleSeacherSumm
    {
        public int SumExists(int sum, string digitsStr)
        {
            var hashMap = new HashSet<int>();
            foreach (var digit in digitsStr.Split(' ').Select(e => int.Parse(e)).ToArray())
            {
                if (hashMap.Contains(digit))
                {
                    return 1;
                }
                hashMap.Add(sum - digit);
            }
            return 0;
        }
    }

    class EconomSeacherSumm
    {
        public int SumExists(int sum, string digitsStr)
        {
            var hashMap = new HashSet<int>();
            int idx;
            var prev = 0;
            while ((idx = digitsStr.IndexOf(" ", prev, StringComparison.Ordinal)) != -1)
            {
                var digit = int.Parse(digitsStr.Substring(prev, idx - prev));
                if (hashMap.Contains(digit))
                {
                    return 1;
                }
                hashMap.Add(sum - digit);
                prev = idx + 1;
            }
            return 0;
        }
    }
}