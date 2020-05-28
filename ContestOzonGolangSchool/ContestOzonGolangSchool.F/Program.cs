using System;
using System.IO;
using System.Linq;

namespace ContestOzonGolangSchool.F
{
    public class Program
    {
        public static void Main(string[] args)
        {
#if DEBUG
            var sw = new System.Diagnostics.Stopwatch();
            sw.Start();
#endif
            using (var inputStream = File.OpenRead("input.txt"))
            using (var inputReader = new StreamReader(inputStream))
            {
                var sum = long.Parse(inputReader.ReadLine());
                var digits = inputReader.ReadLine().Split(' ').Select(e => long.Parse(e)).ToArray();
                var result = FindSum(sum, digits);
                File.WriteAllText("output.txt", result.ToString());
            }
#if DEBUG
            Console.WriteLine("Elapsed: " + sw.ElapsedMilliseconds);
#endif
        }

        private static int FindSum(long sum, long[] digits)
        {
            for (var i = 0; i < digits.Length; i++)
            {
                var first = digits[i];
                foreach (var second in digits.Skip(i + 1))
                {
                    if (first + second == sum)
                    {
                        return 1;
                    }
                }
            }
            return 0;
        }
    }
}
