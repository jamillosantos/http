package hermes

import (
	"sync"
	"time"

	"github.com/lab259/go-rscsrv"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Services", func() {
	Describe("Fasthttp Service", func() {
		It("should not return any error loading the service", func() {
			var service FasthttpService
			result, err := service.LoadConfiguration()
			Expect(err).NotTo(BeNil())
			Expect(err.Error()).To(ContainSubstring("not implemented"))
			Expect(result).To(BeNil())
		})

		It("should apply a given pointer configuration", func() {
			var service FasthttpService
			Expect(service.ApplyConfiguration(&FasthttpServiceConfiguration{
				Bind: "12345",
			})).To(BeNil())
			Expect(service.Configuration.Bind).To(Equal("12345"))
		})

		It("should apply a given configuration", func() {
			var service FasthttpService
			Expect(service.ApplyConfiguration(FasthttpServiceConfiguration{
				Bind: "12345",
			})).To(BeNil())
			Expect(service.Configuration.Bind).To(Equal("12345"))
		})

		It("should fail applying a wrong type configuration", func() {
			var service FasthttpService
			Expect(service.ApplyConfiguration(map[string]interface{}{
				"bind": "12345",
			})).To(Equal(rscsrv.ErrWrongConfigurationInformed))
		})

		It("should start and stop the service", func(done Done) {
			var service FasthttpService
			service.Configuration.Bind = ":32301" // High port
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer GinkgoRecover()

				wg.Done()
				Expect(service.Start()).To(BeNil())
				wg.Done()
			}()
			wg.Wait()
			time.Sleep(time.Millisecond * 500)
			wg.Add(1)
			Expect(service.Stop()).To(BeNil())
			wg.Wait()
			done <- true
		}, 1)

		It("should restart the service that is not started", func(done Done) {
			var service FasthttpService
			service.Configuration.Bind = ":32301" // High port
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer GinkgoRecover()

				wg.Done()
				Expect(service.Restart()).To(BeNil())
				wg.Done()
			}()
			wg.Wait()
			time.Sleep(time.Millisecond * 500)
			wg.Add(1)
			Expect(service.Stop()).To(BeNil())
			wg.Wait()
			done <- true
		}, 1)

		It("should stop a stopped service", func() {
			var service FasthttpService
			Expect(service.Stop()).To(BeNil())
		})

		It("should restart the service", func(done Done) {
			var service FasthttpService
			service.Configuration.Bind = ":32301" // High port

			ch := make(chan string, 10)

			// Just ignore result
			go func() {
				ch <- "service:step1:begin"
				service.Start()
				ch <- "service:step1:end"
			}()

			time.Sleep(time.Millisecond * 500) // Waits for the service a bit

			go func() {
				ch <- "service:step2:begin"
				service.Restart()
				time.Sleep(time.Millisecond * 500)
				ch <- "service:step2:end"
			}()

			time.Sleep(time.Millisecond * 500) // Waits for the service a bit

			Expect(service.Stop()).To(BeNil())
			Expect(<-ch).To(Equal("service:step1:begin"))
			Expect(<-ch).To(Equal("service:step2:begin"))
			Expect(<-ch).To(Equal("service:step1:end"))
			Expect(<-ch).To(Equal("service:step2:end"))
			done <- true
		}, 2)
	})
})
