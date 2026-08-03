export default function AuthLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <main className="flex justify-center bg-background mt-10">
      <div className="w-[60vw] rounded-xl bg-white p-8 shadow-lg">
        {children}
      </div>
    </main>
  );
}
