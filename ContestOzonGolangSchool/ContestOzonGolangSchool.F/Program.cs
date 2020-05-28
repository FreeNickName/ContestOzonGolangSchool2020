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
                var str = inputReader.ReadLine();
#if DEBUG
                Console.WriteLine("Total memory 2: " + GC.GetTotalMemory(false));
#endif
                result = SumExistsEconom(sum, str);
            }
            
            File.WriteAllText("output.txt", result.ToString());

#if DEBUG
            Console.WriteLine("Total memory 3: " + GC.GetTotalMemory(false));
            Console.WriteLine("Elapsed: " + sw.ElapsedMilliseconds);
            Console.WriteLine("Result: " + result);
            Console.ReadLine();
#endif
        }

        private static int SumExistsEconom(int sum, string digitsstr)
        {
            var hashMap = new HashSet<int>();
            int idx;
            //var prev = 0;
            for (var prev = 0; (idx = digitsstr.IndexOf(" ", prev, StringComparison.Ordinal)) != -1; prev = idx +1)
            {
                var digit = int.Parse(digitsstr.Substring(prev, idx - prev));
                if (digit >= sum)
                {
                    continue;
                }
                if (hashMap.Contains(digit))
                {
                    return 1;
                }
                hashMap.Add(sum - digit);
                //prev = idx + 1;
            }
            return 0;
        }

        private static int SumExistsDic(int sum, string digitsStr)
        {
            var hashMap = new Dictionary<int, object>();
            foreach (var digit in digitsStr.Split(' ').Select(e => int.Parse(e)).ToArray())
            {
                if (digit >= sum)
                {
                    continue;
                }
                if (hashMap.ContainsKey(digit))
                {
                    return 1;
                }
                hashMap.Add(sum - digit, null);
            }
            return 0;
        }

        private static int SumExistsHashSet(int sum, int[] digits)
        {
            var hashMap = new HashSet<int>();
            foreach (var digit in digits)
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
