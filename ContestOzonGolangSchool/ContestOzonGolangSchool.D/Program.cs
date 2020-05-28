using System;
using System.IO;
using System.Linq;

namespace ContestOzonGolangSchool.D
{
    class Program
    {
        static void Main(string[] args)
        {
            using (var inputStream = Console.OpenStandardInput())
            using (var inputReader = new StreamReader(inputStream))
            {
                var digits = inputReader.ReadToEnd().Trim().Split(' ').Select(e => e.Select(c => int.Parse(c.ToString())).Reverse().ToArray()).ToArray();
                var digitMin = digits[0].Length > digits[1].Length ? digits[1] : digits[0];
                var digitMax = digits[0].Length > digits[1].Length ? digits[0] : digits[1];
                int[] res = new int[digitMax.Length + 1];
                var overflow = 0;
                for (var i = 0; i < res.Length; i++)
                {
                    if (i >= digitMax.Length)
                    {
                        res[i] = overflow;
                        break;
                    }
                    
                    var sum = (i >= digitMin.Length ? 0 : digitMin[i]) + digitMax[i] + overflow;
                    if (sum > 9)
                    {
                        overflow = 1;
                        res[i] = sum - 10;
                    }
                    else
                    {
                        overflow = 0;
                        res[i] = sum;
                    }
                }
                Console.WriteLine(string.Join("", (res.Last() == 0 ? res.Take(res.Length - 1) : res).Reverse()));
            }
        }
    }
}
