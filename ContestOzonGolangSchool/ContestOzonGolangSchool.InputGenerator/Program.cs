using System;
using System.IO;
using System.Linq;

namespace ContestOzonGolangSchool.InputGenerator
{
    class Program
    {
        static void Main(string[] args)
        {
            var key = args.FirstOrDefault()?.ToLower();
            Console.WriteLine("Key: " + key);
            switch (key)
            {
                case "a":
                    {
                        var res = string.Empty;
                        var random = new Random();
                        var max = int.Parse(args.Skip(1).FirstOrDefault() ?? "10000");
                        for (var i = 0; i < max; i++)
                        {
                            res += random.Next() + Environment.NewLine;
                        }
                        File.WriteAllText("input-201.txt", res);
                        break;
                    }
                case "f":
                    {
                        var random = new Random();
                        var res = random.Next(1, 999999998) + Environment.NewLine;
                        var max = int.Parse(args.Skip(1).FirstOrDefault() ?? "10000");
                        for (var i = 0; i < max; i++)
                        {
                            res += random.Next(1, 999999998) + " ";
                        }
                        File.WriteAllText("input.txt", res.Trim());
                        break;
                    }
                default:
                    {
                        Console.WriteLine("Key unsupported");
                        break;
                    }
            }
        }
    }
}
