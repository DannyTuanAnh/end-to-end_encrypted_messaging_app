import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Link } from "react-router-dom";
import { Smartphone, Zap, Lock, File, Rocket } from "lucide-react";
export default function Landing() {
  const features = [
    {
      title: "Realtime Messaging",
      description:
        "Send and receive messages instantly with stable realtime synchronization.",
      icon: <Zap className="h-6 w-6 text-primary" />,
    },
    {
      title: "Secure Authentication",
      description:
        "Protected sessions using secure HTTP-only cookies and encrypted connections.",
      icon: <Lock className="h-6 w-6 text-primary" />,
    },
    {
      title: "Cross Platform",
      description:
        "Use the application on desktop, tablet, and mobile devices seamlessly.",
      icon: <Smartphone className="h-6 w-6 text-primary" />,
    },
    {
      title: "Media Sharing",
      description:
        "Share images, files, and media with fast upload performance.",
      icon: <File className="h-6 w-6 text-primary" />,
    },
  ];

  const steps = [
    {
      title: "Create Account",
      description: "Sign in securely using your Google account.",
    },
    {
      title: "Connect With Friends",
      description: "Find and start conversations instantly.",
    },
    {
      title: "Chat In Realtime",
      description: "Enjoy fast and smooth realtime messaging.",
    },
  ];

  const testimonials = [
    {
      name: "Sarah Johnson",
      role: "Product Designer",
      comment:
        "The UI is clean, fast, and the realtime updates feel incredibly smooth.",
    },
    {
      name: "Michael Lee",
      role: "Frontend Developer",
      comment:
        "Authentication and session handling work flawlessly across devices.",
    },
    {
      name: "Emily Brown",
      role: "Student",
      comment: "Simple to use and perfect for keeping conversations organized.",
    },
  ];

  return (
    <div className="min-h-screen bg-background text-foreground">
      {/* Header */}
      <header className="sticky top-0 z-50 border-b bg-background/80 backdrop-blur">
        <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
          <div className="flex items-center gap-2">
            <div className="flex h-10 w-10 items-center justify-center rounded-2xl bg-primary text-primary-foreground font-bold">
              C
            </div>
            <div>
              <h1 className="text-lg font-bold">ChatFlow</h1>
              <p className="text-xs text-muted-foreground">
                Modern Messaging Platform
              </p>
            </div>
          </div>

          <nav className="hidden gap-8 md:flex">
            <a href="#features" className="text-sm hover:text-primary">
              Features
            </a>
            <a href="#how-it-works" className="text-sm hover:text-primary">
              How It Works
            </a>
            <a href="#testimonials" className="text-sm hover:text-primary">
              Testimonials
            </a>
            <a href="#faq" className="text-sm hover:text-primary">
              FAQ
            </a>
          </nav>

          <div className="flex items-center gap-3">
            <Link to="/auth/login">
              <Button variant="outline" className="rounded-xl">
                Sign In
              </Button>
            </Link>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 bg-gradient-to-b from-primary/10 to-transparent"></div>

        <div className="relative mx-auto grid max-w-7xl gap-12 px-6 py-20 lg:grid-cols-2 lg:items-center">
          <div>
            <div className="mb-4 inline-flex items-center rounded-full border bg-muted px-4 py-2 text-sm">
              {<Rocket className="h-4 w-4 mr-2 text-primary" />} Fast • Secure •
              Realtime
            </div>

            <h1 className="text-5xl font-bold leading-tight lg:text-6xl">
              Modern Messaging For
              <span className="text-primary"> Everyone</span>
            </h1>

            <p className="mt-6 max-w-xl text-lg text-muted-foreground">
              ChatFlow helps teams, friends, and communities communicate with
              fast realtime messaging, secure authentication, and a beautiful
              modern interface.
            </p>

            <div className="mt-8 flex flex-wrap gap-4">
              <Link to="/auth/login">
                <Button size="lg" className="rounded-2xl shadow-lg">
                  Start Chatting
                </Button>
              </Link>

              <a href="#features">
                <Button variant="outline" size="lg" className="rounded-2xl">
                  Learn More
                </Button>
              </a>
            </div>

            <div className="mt-10 flex flex-wrap gap-8 text-sm text-muted-foreground">
              <div>
                <p className="text-2xl font-bold text-foreground">10K+</p>
                <p>Active Users</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-foreground">99.9%</p>
                <p>Uptime</p>
              </div>
              <div>
                <p className="text-2xl font-bold text-foreground">24/7</p>
                <p>Realtime Delivery</p>
              </div>
            </div>
          </div>

          {/* Hero Chat Preview */}
          <div className="relative">
            <Card className="rounded-3xl shadow-2xl">
              <CardContent className="p-6">
                <div className="mb-6 flex items-center justify-between border-b pb-4">
                  <div className="flex items-center gap-3">
                    <div className="h-12 w-12 rounded-full bg-primary/20" />
                    <div>
                      <h3 className="font-semibold">Development Team</h3>
                      <p className="text-sm text-green-500">Online</p>
                    </div>
                  </div>

                  <div className="flex gap-2">
                    <div className="h-3 w-3 rounded-full bg-red-400" />
                    <div className="h-3 w-3 rounded-full bg-yellow-400" />
                    <div className="h-3 w-3 rounded-full bg-green-400" />
                  </div>
                </div>

                <div className="space-y-4">
                  <div className="max-w-xs rounded-2xl bg-muted p-4">
                    <p className="text-sm">
                      Hey team 👋 The new frontend deployment is ready.
                    </p>
                  </div>

                  <div className="ml-auto max-w-xs rounded-2xl bg-primary p-4 text-primary-foreground">
                    <p className="text-sm">
                      Great! I’ll connect the auth flow with the backend now.
                    </p>
                  </div>

                  <div className="max-w-xs rounded-2xl bg-muted p-4">
                    <p className="text-sm">
                      Realtime notifications are already working perfectly ⚡
                    </p>
                  </div>
                </div>

                <div className="mt-6 flex items-center gap-3 border-t pt-4">
                  <Input
                    placeholder="Type a message..."
                    className="flex-1 rounded-xl"
                  />

                  <Button className="rounded-xl">Send</Button>
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="mx-auto max-w-7xl px-6 py-20">
        <div className="mb-14 text-center">
          <h2 className="text-4xl font-bold">Powerful Features</h2>
          <p className="mt-4 text-muted-foreground">
            Everything you need for modern communication.
          </p>
        </div>

        <div className="grid gap-6 md:grid-cols-2 xl:grid-cols-4">
          {features.map((feature) => (
            <Card
              key={feature.title}
              className="rounded-3xl shadow-sm transition hover:-translate-y-1 hover:shadow-lg"
            >
              <CardContent className="p-6">
                <div className="mb-4 text-4xl">{feature.icon}</div>

                <h3 className="text-xl font-semibold">{feature.title}</h3>

                <p className="mt-3 text-sm leading-6 text-muted-foreground">
                  {feature.description}
                </p>
              </CardContent>
            </Card>
          ))}
        </div>
      </section>

      {/* How It Works */}
      <section id="how-it-works" className="bg-muted/40 px-6 py-20">
        <div className="mx-auto max-w-6xl">
          <div className="mb-14 text-center">
            <h2 className="text-4xl font-bold">How It Works</h2>
            <p className="mt-4 text-muted-foreground">
              Start chatting in just a few simple steps.
            </p>
          </div>

          <div className="grid gap-8 md:grid-cols-3">
            {steps.map((step, index) => (
              <div
                key={step.title}
                className="rounded-3xl border bg-background p-8 text-center shadow-sm"
              >
                <div className="mx-auto mb-6 flex h-14 w-14 items-center justify-center rounded-full bg-primary text-lg font-bold text-primary-foreground">
                  {index + 1}
                </div>

                <h3 className="text-xl font-semibold">{step.title}</h3>

                <p className="mt-3 text-muted-foreground">{step.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Testimonials */}
      <section id="testimonials" className="mx-auto max-w-7xl px-6 py-20">
        <div className="mb-14 text-center">
          <h2 className="text-4xl font-bold">What Users Say</h2>
          <p className="mt-4 text-muted-foreground">
            Trusted by developers, students, and teams.
          </p>
        </div>

        <div className="grid gap-6 lg:grid-cols-3">
          {testimonials.map((testimonial) => (
            <div
              key={testimonial.name}
              className="rounded-3xl border bg-card p-8 shadow-sm"
            >
              <p className="leading-7 text-muted-foreground">
                “{testimonial.comment}”
              </p>

              <div className="mt-6 flex items-center gap-4">
                <div className="h-12 w-12 rounded-full bg-primary/20" />

                <div>
                  <h4 className="font-semibold">{testimonial.name}</h4>
                  <p className="text-sm text-muted-foreground">
                    {testimonial.role}
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>
      </section>

      {/* CTA */}
      <section className=" border bg-primary px-8 py-16 text-center text-primary-foreground">
        <h2 className="text-4xl font-bold">Ready To Start Messaging?</h2>

        <p className="mx-auto mt-4 max-w-2xl text-primary-foreground/80">
          Create your account and experience fast, secure, and modern realtime
          communication today.
        </p>

        <div className="mt-8 flex flex-wrap justify-center gap-4">
          <Link to="/auth/login">
            <button className="rounded-2xl bg-background px-6 py-3 font-medium text-foreground hover:opacity-90 transition">
              Get Started
            </button>
          </Link>

          <Link to="/contact">
            <button className="rounded-2xl border border-primary-foreground/30 px-6 py-3 font-medium hover:bg-primary-foreground/10 transition">
              Contact Us
            </button>
          </Link>
        </div>
      </section>

      {/* FAQ */}
      <section id="faq" className="bg-muted/40 px-6 py-20">
        <div className="mx-auto max-w-4xl">
          <div className="mb-14 text-center">
            <h2 className="text-4xl font-bold">Frequently Asked Questions</h2>
          </div>

          <div className="space-y-4">
            {[
              {
                q: "Is ChatFlow free to use?",
                a: "Yes, the core messaging features are available for free.",
              },
              {
                q: "Does it support realtime notifications?",
                a: "Yes, realtime updates are powered by SSE technology.",
              },
              {
                q: "Is authentication secure?",
                a: "Authentication uses secure HTTP-only cookies and encrypted HTTPS connections.",
              },
            ].map((item) => (
              <div
                key={item.q}
                className="rounded-2xl border bg-background p-6"
              >
                <h3 className="text-lg font-semibold">{item.q}</h3>
                <p className="mt-2 text-muted-foreground">{item.a}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t px-6 py-10">
        <div className="mx-auto flex max-w-7xl flex-col gap-6 md:flex-row md:items-center md:justify-between">
          <div>
            <h3 className="text-xl font-bold">ChatFlow</h3>
            <p className="mt-2 text-sm text-muted-foreground">
              Modern realtime messaging platform.
            </p>
          </div>

          <div className="flex flex-wrap gap-6 text-sm text-muted-foreground">
            <a href="#">Privacy Policy</a>
            <a href="#">Terms of Service</a>
            <a href="#">Support</a>
            <a href="#">Contact</a>
          </div>
        </div>
      </footer>
    </div>
  );
}
