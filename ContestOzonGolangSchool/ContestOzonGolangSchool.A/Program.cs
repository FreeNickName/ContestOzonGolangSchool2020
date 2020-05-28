using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;

namespace ContestOzonGolangSchool.A
{
    public class Program
    {
        public static void Main(string[] args)
        {
#if DEBUG
            var sw = new System.Diagnostics.Stopwatch();
            sw.Start();
#endif
            var input = File.ReadAllText("input-201.txt");
            var digits = input.Split(Environment.NewLine, StringSplitOptions.RemoveEmptyEntries).Select(e => int.Parse(e));
            var successPairs = new HashSet<int>();
            foreach (var digit in digits)
            {
                if (successPairs.Contains(digit))
                {
                    successPairs.Remove(digit);
                }
                else
                {
                    successPairs.Add(digit);
                }
            }
#if DEBUG
            Console.WriteLine(sw.ElapsedMilliseconds);
#endif
            Console.Write(successPairs.Count > 0 ? successPairs.First().ToString() : "None");
        }
    }
}
