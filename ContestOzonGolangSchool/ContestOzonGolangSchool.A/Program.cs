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
            var incompletePairs = new HashSet<int>();
            using (var inputStream = Console.OpenStandardInput())
            using (var inputReader = new StreamReader(inputStream))
            {
                var line = string.Empty;
                while ((line = inputReader.ReadLine()) != null)
                {
                    var digit = int.Parse(line);
                    if (incompletePairs.Contains(digit))
                    {
                        incompletePairs.Remove(digit);
                    }
                    else
                    {
                        incompletePairs.Add(digit);
                    }
                }
            }
#if DEBUG
            Console.WriteLine(sw.ElapsedMilliseconds);
#endif
            Console.WriteLine(incompletePairs.Count > 0 ? incompletePairs.First().ToString() : "None");
        }
    }
}
