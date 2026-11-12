package closer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/DoomGuy1818/rocket-science/platform/pkg/logger"
)

const shutdownTimeout = 5 * time.Second

type Logger interface {
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
}

type Closer struct {
	mu     sync.Mutex
	once   sync.Once
	funcs  []func(context.Context) error
	done   chan struct{}
	logger Logger
}

var globalCloser = NewWithLogger(&logger.NoopLogger{})

// Configure настраивает глобальный closer для обработки системных сигналов
func Configure(signals ...os.Signal) {
	go globalCloser.handleSignals(signals...)
}

func New(signals ...os.Signal) *Closer {
	return NewWithLogger(logger.Logger(), signals...)
}

func NewWithLogger(logger Logger, signals ...os.Signal) *Closer {
	c := &Closer{
		done:   make(chan struct{}),
		logger: logger,
	}

	if len(signals) > 0 {
		go c.handleSignals(signals...)
	}

	return c
}

// AddNamed добавляет функцию закрытия с именем зависимости для логирования в глобальный closer
func (c *Closer) AddNamed(name string, f func(context.Context) error) {
	c.Add(
		func(ctx context.Context) error {
			startTime := time.Now()

			c.logger.Info(ctx, fmt.Sprintf("🧩 Закрываем %s...", name))

			err := f(ctx)

			duration := time.Since(startTime)
			if err != nil {
				c.logger.Error(ctx, fmt.Sprintf("❌ Ошибка при закрытии %s: %v (заняло %s)", name, err, duration))
			} else {
				c.logger.Info(ctx, fmt.Sprintf("✅ %s успешно закрыт за %s", name, duration))
			}

			return err
		},
	)
}

// Add добавляет функции закрытия в глобальный closer
func (c *Closer) Add(f ...func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.funcs = append(c.funcs, f...)
}

// CloseAll вызывает все зарегистрированные функции закрытия.
// Возвращает первую возникшую ошибку, если таковая была.
func (c *Closer) CloseAll(ctx context.Context) error {
	var result error

	c.once.Do(
		func() {
			defer close(c.done)

			c.mu.Lock()
			funcs := c.funcs
			c.funcs = nil
			defer c.mu.Unlock()

			if len(funcs) == 0 {
				c.logger.Info(ctx, "ℹ️ Нет функций для закрытия.")
				return
			}

			c.logger.Info(ctx, "🚦 Начинаем процесс graceful shutdown...")
			errCh := make(chan error, len(funcs))

			var wg sync.WaitGroup

			for i := len(funcs) - 1; i >= 0; i-- {
				f := funcs[i]

				wg.Add(1)

				go func(f func(context.Context) error) {
					defer wg.Done()

					defer func() {
						if r := recover(); r != nil {
							errCh <- errors.New("panic recovered in closer")
							c.logger.Error(ctx, "⚠️ Panic в функции закрытия", zap.Any("error", r))
						}
					}()

					if err := f(ctx); err != nil {
						errCh <- err
					}
				}(f)
			}

			go func() {
				wg.Wait()
				defer close(errCh)
			}()

			select {
			case err, ok := <-errCh:
				if !ok {
					c.logger.Info(ctx, "✅ Все ресурсы успешно закрыты")
					return
				}
				c.logger.Error(ctx, "❌ Ошибка при закрытии", zap.Error(err))
				if result == nil {
					result = err
				}
			case <-ctx.Done():
				c.logger.Info(ctx, "⚠️ Контекст отменён во время закрытия", zap.Error(ctx.Err()))
				if result == nil {
					result = ctx.Err()
				}
				return
			}
		},
	)

	return result
}

func AddNamed(name string, f func(context.Context) error) {
	globalCloser.AddNamed(name, f)
}

func SetLogger(l Logger) {
	globalCloser.SetLogger(l)
}

// Add добавляет функции закрытия в глобальный closer
func Add(f ...func(context.Context) error) {
	globalCloser.Add(f...)
}

// CloseAll инициирует процесс закрытия всех зарегистрированных функций глобального closer'а
func CloseAll(ctx context.Context) error {
	return globalCloser.CloseAll(ctx)
}

func (c *Closer) SetLogger(l Logger) {
	c.logger = l
}

func (c *Closer) handleSignals(signals ...os.Signal) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	defer signal.Stop(ch)

	select {
	case <-ch:
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer shutdownCancel()

		c.logger.Info(context.Background(), "🛑 Получен системный сигнал, начинаем graceful shutdown...")

		if err := c.CloseAll(shutdownCtx); err != nil {
			c.logger.Error(context.Background(), "❌ Ошибка при закрытии ресурсов: %v", zap.Error(err))
		}

	case <-c.done:
	}
}
