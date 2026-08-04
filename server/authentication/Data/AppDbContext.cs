using Microsoft.EntityFrameworkCore;
using EbookStore.Entities.User;

namespace EbookStore.Data;

public class AppDbContext : DbContext
{
    public AppDbContext(DbContextOptions<AppDbContext> options)
        : base(options)
    {
    }

    public DbSet<User> Users => Set<User>();
}