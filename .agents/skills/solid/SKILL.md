---
name: solid
description: Use this skill when writing code, implementing features, refactoring, planning architecture, designing systems, reviewing code, or debugging. This skill transforms junior-level code into senior-engineer quality software through SOLID principles, TDD, clean code practices, and professional software design.
---

# Solid Skills: Professional Software Engineering (PHP 8.3 / Laravel 12)

You are now operating as a senior software engineer. Every line of code you write, every design decision you make, and every refactoring you perform must embody professional craftsmanship.

## When This Skill Applies

**ALWAYS use this skill when:**
- Writing ANY code (features, fixes, utilities)
- Refactoring existing code
- Planning or designing architecture
- Reviewing code quality
- Debugging issues
- Creating tests
- Making design decisions

## Core Philosophy

> "Code is to create products for users & customers. Testable, flexible, and maintainable code that serves the needs of the users is GOOD because it can be cost-effectively maintained by developers."

The goal of software: Enable developers to **discover, understand, add, change, remove, test, debug, deploy**, and **monitor** features efficiently.

## The Non-Negotiable Process

### 1. ALWAYS Start with Tests (TDD)

**Red-Green-Refactor is not optional:**

```
1. RED    - Write a failing test that describes the behavior
2. GREEN  - Write the SIMPLEST code to make it pass
3. REFACTOR - Clean up, remove duplication (Rule of Three)
```

**The Three Laws of TDD:**
1. You cannot write production code unless it makes a failing test pass
2. You cannot write more test code than is sufficient to fail
3. You cannot write more production code than is sufficient to pass

**Design happens during REFACTORING, not during coding.**

**Pest example (project standard):**
```php
it('applies 10% discount when order exceeds 100', function () {
    // Arrange
    $service = new CalculateDiscountAction();
    $dto = new OrderDTO(amount: 150.0);

    // Act
    $result = $service->handle($dto);

    // Assert
    expect($result)->toBe(135.0);
});
```

### 2. Apply SOLID Principles Rigorously

Every class, every module, every function:

| Principle | Question to Ask | Laravel Example |
|-----------|-----------------|-----------------|
| **S**RP - Single Responsibility | "Does this have ONE reason to change?" | One Action = one use case |
| **O**CP - Open/Closed | "Can I extend without modifying?" | Strategy pattern via interfaces |
| **L**SP - Liskov Substitution | "Can subtypes replace base types safely?" | Repository interfaces |
| **I**SP - Interface Segregation | "Are clients forced to depend on unused methods?" | Granular interfaces per entity |
| **D**IP - Dependency Inversion | "Do high-level modules depend on abstractions?" | Bind interfaces in ServiceProvider |

**SRP in Laravel — Action pattern:**
```php
<?php

declare(strict_types=1);

namespace App\Actions\Order;

use App\DTOs\Order\CreateOrderDTO;
use App\Interfaces\Repositories\OrderRepositoryInterface;
use App\Models\Order;
use Lorisleiva\Actions\Concerns\AsAction;

// ✅ Single responsibility: create an order. Nothing else.
class CreateOrderAction
{
    use AsAction;

    public function __construct(
        private readonly OrderRepositoryInterface $repository,
    ) {}

    public function handle(CreateOrderDTO $dto): Order
    {
        return $this->repository->create($dto->toArray());
    }
}
```

**OCP — Open for extension via interfaces:**
```php
<?php

declare(strict_types=1);

namespace App\Interfaces;

// ✅ New payment methods extend without modifying existing code
interface PaymentGatewayInterface
{
    public function charge(Money $amount, PaymentMethod $method): PaymentResult;
}

// Each gateway implements the interface — no switch/if-else chains
class StripeGateway implements PaymentGatewayInterface { /* ... */ }
class PayPalGateway implements PaymentGatewayInterface { /* ... */ }
```

**DIP — Dependency Inversion in Laravel:**
```php
// AppServiceProvider.php — bind abstractions to concretions
$this->app->bind(OrderRepositoryInterface::class, OrderRepository::class);
$this->app->bind(PaymentGatewayInterface::class, StripeGateway::class);
```

### 3. Write Clean, Human-Readable Code

**Naming (in order of priority):**
1. **Consistency** - Same concept = same name everywhere
2. **Understandability** - Domain language, not technical jargon
3. **Specificity** - Precise, not vague (avoid `$data`, `$info`, `$manager`)
4. **Brevity** - Short but not cryptic
5. **Searchability** - Unique, greppable names

**Structure:**
- One level of indentation per method
- No `else` keyword when possible (early returns / guard clauses)
- Use strict comparisons (`===`, `!==`) — never rely on loose truthiness
- **ALWAYS wrap primitives in Value Objects or Enums** for domain concepts
- First-class collections (wrap arrays in dedicated collection classes)
- Keep entities small (< 50 lines for classes, < 10 for methods)
- Prefer `DateTimeImmutable` over mutable `DateTime`
- Use constructor property promotion with `readonly`

**Value Objects are MANDATORY for domain concepts (PHP 8.3):**
```php
<?php

declare(strict_types=1);

namespace App\ValueObjects;

use InvalidArgumentException;

// ✅ Value Object — immutable, self-validating, type-safe
final readonly class Email
{
    public function __construct(
        public string $value,
    ) {
        if (!filter_var($value, FILTER_VALIDATE_EMAIL)) {
            throw new InvalidArgumentException("Invalid email: {$value}");
        }
    }

    public function equals(self $other): bool
    {
        return $this->value === $other->value;
    }

    public function __toString(): string
    {
        return $this->value;
    }
}

// ✅ Value Object for monetary amounts
final readonly class Money
{
    public function __construct(
        public int $cents,
        public string $currency = 'COP',
    ) {
        if ($this->cents < 0) {
            throw new InvalidArgumentException('Amount cannot be negative');
        }
    }

    public function add(self $other): self
    {
        if ($this->currency !== $other->currency) {
            throw new InvalidArgumentException('Cannot add different currencies');
        }

        return new self($this->cents + $other->cents, $this->currency);
    }
}
```

**Enums replace magic strings and constants (PHP 8.1+):**
```php
<?php

declare(strict_types=1);

namespace App\Enums;

// ✅ Backed enum — type-safe, serializable, exhaustive
enum OrderStatus: string
{
    case Pending = 'pending';
    case Approved = 'approved';
    case Rejected = 'rejected';
    case Delivered = 'delivered';

    public function label(): string
    {
        return match ($this) {
            self::Pending => 'Pendiente',
            self::Approved => 'Aprobado',
            self::Rejected => 'Rechazado',
            self::Delivered => 'Entregado',
        };
    }

    public function canTransitionTo(self $target): bool
    {
        return match ($this) {
            self::Pending => in_array($target, [self::Approved, self::Rejected]),
            self::Approved => $target === self::Delivered,
            default => false,
        };
    }
}

// NEVER use raw strings for statuses:
// ❌ $order->status = 'approved';
// ✅ $order->status = OrderStatus::Approved;
```

### 4. Design with Responsibility in Mind

**Ask these questions for every class:**
1. "What pattern is this?" (Action, Service, Repository, DTO, Observer, etc.)
2. "Is it doing too much?" (Check method count and line count)

**Object Stereotypes in Laravel:**
- **DTO (Data)** — Holds validated input, immutable (`Spatie\LaravelData\Data`)
- **Action** — Single business operation, reusable (`lorisleiva/laravel-actions`)
- **Service** — Complex workflow coordinating multiple Actions/entities
- **Repository** — DB access only, no business logic
- **Resource** — API output transformation
- **Observer** — Reacts to model lifecycle events (cache, events)
- **Job** — Thin wrapper that calls an Action asynchronously

**Laravel-specific responsibility boundaries:**
```php
// ❌ Controller doing too much (God controller)
public function store(Request $request): JsonResponse
{
    $validated = $request->validate([...]);
    $order = Order::create($validated);
    $order->items()->createMany($validated['items']);
    Cache::forget('orders');
    Mail::to($order->user)->send(new OrderCreated($order));
    return response()->json($order);
}

// ✅ Each layer has ONE job
public function store(StoreOrderRequest $request): JsonResponse
{
    $order = CreateOrderAction::run(
        CreateOrderDTO::from($request->validated())
    );

    return (new OrderResource($order))->response()->setStatusCode(201);
}
// Validation → FormRequest
// Input structure → DTO
// Business logic → Action
// DB access → Repository (inside Action)
// Cache invalidation → Observer
// Email → Queued notification (inside Action or Observer)
// Output → Resource
```

### 5. Manage Complexity Ruthlessly

**Essential complexity** = inherent to the problem domain
**Accidental complexity** = introduced by our solutions

**Detect complexity through:**
- Change amplification (small change = many files)
- Cognitive load (hard to understand)
- Unknown unknowns (surprises in behavior)

**Fight complexity with:**
- YAGNI - Don't build what you don't need NOW
- KISS - Simplest solution that works
- DRY - But only after Rule of Three (wait for 3 duplications)

**Laravel-specific complexity traps:**
```php
// ❌ Over-engineering: creating a Repository for a query used once
class OrderStatsQueryRepository
{
    public function getMonthlyTotal(): int { /* used in one place */ }
}

// ✅ Keep it in the Action until you need it elsewhere
class GetMonthlyOrderStatsAction
{
    use AsAction;

    public function handle(): int
    {
        return Order::where('created_at', '>=', now()->startOfMonth())->count();
    }
}

// ❌ Premature abstraction: interface for a class with one implementation
interface OrderCalculatorInterface { }
class OrderCalculator implements OrderCalculatorInterface { }

// ✅ Add the interface when you actually have 2+ implementations
class OrderCalculator { }
```

### 6. Architect for Change

**Vertical Slicing:**
- Features as end-to-end slices (DTO → Action → Repository → Resource)
- Each feature self-contained in its entity folder

**Horizontal Decoupling:**
- Layers don't know about each other's internals
- Dependencies point inward (toward domain)

**The Dependency Rule in Laravel:**
```
Controller → Action → Repository → Model
     ↓           ↓          ↓
FormRequest    DTO     Interface (bound in ServiceProvider)
     ↓
  Resource (output only)
```

- Controllers depend on Actions (never on Repositories directly)
- Actions depend on Repository interfaces (never on Eloquent directly in complex cases)
- Models are the innermost layer — they know nothing about HTTP

## The Four Elements of Simple Design (XP)

In priority order:
1. **Runs all the tests** - `php vendor/bin/pest` must pass
2. **Expresses intent** - Readable, reveals purpose
3. **No duplication** - DRY (but Rule of Three)
4. **Minimal** - Fewest classes, methods possible

## Code Smell Detection

**Stop and refactor when you see:**

| Smell | Solution | Laravel Example |
|-------|----------|-----------------|
| Long Method | Extract methods, compose method pattern | Action `handle()` > 25 lines → split |
| Large Class | Extract class, single responsibility | Controller with business logic → extract Action |
| Long Parameter List | Introduce parameter object | 4+ params → create a DTO |
| Divergent Change | Split into focused classes | Model with query scopes + business logic → extract QueryRepository |
| Shotgun Surgery | Move related code together | Cache invalidation scattered → centralize in Observer |
| Feature Envy | Move method to the envied class | Controller accessing model internals → move to Action |
| Data Clumps | Extract class for grouped data | `$lat, $lng, $radius` always together → `LocationDTO` |
| Primitive Obsession | Wrap in Value Objects or Enums | `string $status` → `OrderStatus` enum |
| Switch Statements | Replace with polymorphism or `match` | `if/elseif` chain on type → Strategy + interface |
| Speculative Generality | YAGNI - remove unused abstractions | Interface with one implementation → remove interface |

## Design Patterns in PHP/Laravel

**Creational:** Factory (Eloquent factories, `App::make()`), Builder (Query Builder, Spatie Data)
**Structural:** Adapter (Payment gateways), Decorator (Middleware), Composite (Validation rules)
**Behavioral:** Strategy (interface + multiple implementations), Observer (Model observers), Template Method (base Action classes)

**Warning:** Don't force patterns. Let them emerge from refactoring.

**Laravel-native patterns you already use:**
```php
// Strategy — via interface binding
$this->app->bind(NotificationChannelInterface::class, function () {
    return match (config('notifications.default')) {
        'slack' => new SlackChannel(),
        'email' => new EmailChannel(),
        default => new DatabaseChannel(),
    };
});

// Observer — model lifecycle hooks
#[ObservedBy(OrderObserver::class)]
class Order extends Model { }

// Builder — Spatie Query Builder for complex queries
QueryBuilder::for(Order::class)
    ->allowedFilters(['status', 'user_id'])
    ->allowedSorts(['created_at', 'total'])
    ->paginate();

// Template Method — base Action with shared behavior
abstract class BaseExportAction
{
    use AsAction;

    abstract protected function query(): Builder;
    abstract protected function transform(Model $model): array;

    public function handle(): Collection
    {
        return $this->query()->cursor()->map(fn ($m) => $this->transform($m));
    }
}
```

## Testing Strategy (Pest 4 + Laravel 12)

**Test Types (from inner to outer):**
1. **Unit Tests** - Single class/function, fast, isolated (mock dependencies)
2. **Integration/Feature Tests** - Full HTTP stack, real DB (SQLite in-memory)

**Arrange-Act-Assert Pattern (Pest syntax):**
```php
it('creates an order and returns the model', function () {
    // Arrange — prepare context
    $dto = new CreateOrderDTO(product_id: 1, quantity: 3, amount: 150_00);
    $expectedOrder = Order::factory()->make(['amount' => 150_00]);

    $repo = Mockery::mock(OrderRepositoryInterface::class);
    $repo->shouldReceive('create')->once()->andReturn($expectedOrder);

    $action = new CreateOrderAction($repo);

    // Act — execute the behavior under test
    $result = $action->handle($dto);

    // Assert — verify the observable outcome
    expect($result->amount)->toBe(150_00);
});
```

**Test Naming:** Use concrete examples, not abstract statements
```php
// ❌ BAD: vague, doesn't describe the scenario
it('works correctly', ...);
it('test order', ...);

// ✅ GOOD: describes behavior + condition
it('returns zero discount when order total is below minimum', ...);
it('throws InvalidArgumentException when email format is invalid', ...);
it('dispatches notification job when order status changes to approved', ...);
```

**What to mock vs what to hit:**
```php
// ✅ Mock: external dependencies (I/O boundaries)
Http::fake(['api.external.com/*' => Http::response(['ok' => true])]);
Storage::fake('s3');
Mail::fake();
Queue::fake();
Notification::fake();

// ✅ Real DB: only in Repository tests and Feature tests
uses(Tests\TestCase::class, RefreshDatabase::class);

// ✅ Action::fake(): in Job tests
ProcessOrderAction::fake();
(new ProcessOrderJob(orderId: 42))->handle();
ProcessOrderAction::assertRan();
```

## Behavioral Principles

- **Tell, Don't Ask** - Command objects, don't query and decide externally
- **Design by Contract** - Preconditions (DTO validation), postconditions (return types), invariants (enums)
- **Hollywood Principle** - "Don't call us, we'll call you" (Laravel's IoC container, Events, Observers)
- **Law of Demeter** - Only talk to immediate friends (no `$order->user->address->city->name`)

```php
// ❌ Violates Law of Demeter (train wreck)
$cityName = $order->user->address->city->name;

// ✅ Ask the object directly
$cityName = $order->getShippingCityName();

// ❌ Ask then decide (violates Tell, Don't Ask)
if ($order->getStatus() === OrderStatus::Pending) {
    $order->setStatus(OrderStatus::Approved);
    $order->setApprovedAt(now());
}

// ✅ Tell the object what to do
$order->approve(); // encapsulates the state transition + invariants
```

## Pre-Code Checklist

Before writing ANY code, answer:

1. [ ] Do I understand the requirement? (Write acceptance criteria first)
2. [ ] What test will I write first?
3. [ ] What is the simplest solution?
4. [ ] What patterns might apply? (Don't force them)
5. [ ] Am I solving a real problem or a hypothetical one?
6. [ ] Does this need a transaction? (Multiple writes = yes)
7. [ ] Does this need an atomic lock? (Concurrent access to same resource = yes)
8. [ ] Should any part be queued? (Emails, external APIs, heavy processing = yes)

## During-Code Checklist

While coding, continuously ask:

1. [ ] Is this the simplest thing that could work?
2. [ ] Does this class have a single responsibility?
3. [ ] Am I depending on abstractions or concretions?
4. [ ] Can I name this more clearly?
5. [ ] Is there duplication I should extract? (Rule of Three)
6. [ ] Am I using `readonly` properties where possible?
7. [ ] Am I using enums instead of magic strings?
8. [ ] Does every method have explicit return types?

## Post-Code Checklist

After the code works:

1. [ ] Do all tests pass? (`php vendor/bin/pest --no-coverage`)
2. [ ] Is there any dead code to remove?
3. [ ] Can I simplify any complex conditions with `match` or early returns?
4. [ ] Are names still accurate after changes?
5. [ ] Would a junior understand this in 6 months?
6. [ ] Is `declare(strict_types=1)` at the top of every new file?
7. [ ] Are all new interfaces bound in a ServiceProvider?

## Red Flags - Stop and Rethink

- Writing code without a test
- Controller method longer than 10 lines
- Action `handle()` longer than 25 lines
- More than one level of indentation in a method
- Using `else` when early return works
- Hardcoding values that should be in config or enums
- Creating abstractions before the third duplication
- Adding features "just in case" (YAGNI)
- Depending on concrete implementations instead of interfaces
- God classes that know everything
- Passing `$request` directly to an Action (use DTOs)
- Catching `Throwable` in application code (let it bubble to the handler)
- Using `->get()` on unbounded queries (use `chunk()`, `cursor()`, or `lazy()`)
- Missing `$tries`, `$timeout`, `$backoff` on a Job

## PHP 8.3 Features to Leverage

```php
// ✅ Readonly classes — entire class is immutable
final readonly class Coordinates
{
    public function __construct(
        public float $latitude,
        public float $longitude,
    ) {}
}

// ✅ Typed class constants (PHP 8.3)
class CacheConfig
{
    public const int DEFAULT_TTL = 3600;
    public const string PREFIX = 'app:';
}

// ✅ #[Override] attribute — catch broken parent contracts at compile time
class OrderRepository implements OrderRepositoryInterface
{
    #[Override]
    public function create(array $data): Order
    {
        return Order::create($data);
    }
}

// ✅ match expression — exhaustive, no fall-through
$label = match ($status) {
    OrderStatus::Pending => 'Pendiente',
    OrderStatus::Approved => 'Aprobado',
    OrderStatus::Rejected => 'Rechazado',
    OrderStatus::Delivered => 'Entregado',
};

// ✅ Named arguments — clarity in constructor calls
$dto = new CreateOrderDTO(
    productId: $validated['product_id'],
    quantity: $validated['quantity'],
    notes: $validated['notes'] ?? null,
);

// ✅ Null-safe operator — avoid nested null checks
$cityName = $order->user?->address?->city?->name;

// ✅ First-class callable syntax
$names = $users->map($this->formatName(...));
```

## Remember

> "A little bit of duplication is 10x better than the wrong abstraction."

> "Focus on WHAT needs to happen, not HOW it needs to happen."

> "Design principles become second nature through practice. Eventually, you won't think about SOLID — you'll just write SOLID code."

The journey: Code-first → Best-practice-first → Pattern-first → Responsibility-first → **Systems Thinking**

Your goal is to reach systems thinking — where principles are internalized and you focus on optimizing the entire development process.
