using System;
using System.IO;
using System.Linq;

namespace SampleContest.A
{
    class Program
    {
        static void Main(string[] args)
        {
            var input = File.ReadAllText("input.txt");
            var sum = input.Trim().Split(' ').Select(e => int.Parse(e)).Sum();
            Console.Write(sum);
        }
    }
}
